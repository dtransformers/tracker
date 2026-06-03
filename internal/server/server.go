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

	server := &http.Server{
		Addr:              net.JoinHostPort(s.Host, fmt.Sprintf("%d", s.Port)),
		Handler:           handler,
		ReadTimeout:       defaultIfZero(s.ReadTimeout, 15*time.Second),
		ReadHeaderTimeout: defaultIfZero(s.ReadHeaderTimeout, 10*time.Second),
		WriteTimeout:      defaultIfZero(s.WriteTimeout, 30*time.Second),
		IdleTimeout:       defaultIfZero(s.IdleTimeout, 2*time.Minute),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	if s.TLS.Enabled {
		if s.TLS.CertFile == "" || s.TLS.KeyFile == "" {
			return errors.New("TLS enabled but cert/key file not provided")
		}
		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	errCh := make(chan error, 1)
	go func() {
		if s.TLS.Enabled {
			errCh <- server.ListenAndServeTLS(s.TLS.CertFile, s.TLS.KeyFile)
			return
		}
		errCh <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-stop:
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, defaultIfZero(s.ShutdownTimeout, 10*time.Second))
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
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
