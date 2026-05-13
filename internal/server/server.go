package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TLSSettings struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type Settings struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	TLS               TLSSettings
}

// Run starts the server and blocks until shutdown.
// If TLS is enabled, the server will negotiate HTTP/2 automatically.
func Run(ctx context.Context, handler http.Handler, s Settings) error {
	if handler == nil {
		return errors.New("handler is nil")
	}

	addr := net.JoinHostPort(s.Host, fmt.Sprintf("%d", s.Port))

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       defaultIfZero(s.ReadTimeout, 15*time.Second),
		ReadHeaderTimeout: defaultIfZero(s.ReadHeaderTimeout, 10*time.Second),
		WriteTimeout:      defaultIfZero(s.WriteTimeout, 30*time.Second),
		IdleTimeout:       defaultIfZero(s.IdleTimeout, 2*time.Minute),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	// TLS -> HTTP/2 is enabled by default in Go.
	// We set a minimal TLS config; Go will add "h2" automatically when appropriate.
	if s.TLS.Enabled {
		if s.TLS.CertFile == "" || s.TLS.KeyFile == "" {
			return errors.New("TLS enabled but cert/key file not provided")
		}
		srv.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	// Graceful shutdown on SIGINT/SIGTERM
	errCh := make(chan error, 1)
	go func() {
		if s.TLS.Enabled {
			errCh <- srv.ListenAndServeTLS(s.TLS.CertFile, s.TLS.KeyFile)
			return
		}
		errCh <- srv.ListenAndServe()
	}()

	// Wait for stop signal or server error
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		// external cancellation
	case <-stop:
		// OS signal
	case err := <-errCh:
		// server exited early
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultIfZero(s.ShutdownTimeout, 10*time.Second))
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	// If the server goroutine returns after Shutdown, consume its error
	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server after shutdown: %w", err)
		}
	default:
	}

	return nil
}

func defaultIfZero[T comparable](v, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}
