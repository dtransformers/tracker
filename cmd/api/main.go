package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"tracking/internal/config"
	"tracking/internal/server"

	"github.com/arl/statsviz"
)

func startApp(host string, port int, enabled bool) {
	log.Printf("🚀 starting server on %s:%d (tls=%v)",
		host, port, enabled,
	)
}

func setupVisualizer(server *http.ServeMux, host string, port int) {
	fmt.Printf("🚀 starting statsviz server, Check http://%s:%d/debug/statsviz/", host, port)
	statsviz.Register(server)
}

func runServer(cfg config.Config) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	s := server.Settings{
		Host:              cfg.App.Host,
		Port:              cfg.App.Port,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
		ShutdownTimeout:   10 * time.Second,

		TLS: server.TLSSettings{
			Enabled:  cfg.TLS.Enabled,
			CertFile: cfg.TLS.CertFile,
			KeyFile:  cfg.TLS.KeyFile,
		},
	}
	startApp(cfg.App.Host, cfg.App.Port, cfg.TLS.Enabled)
	if cfg.App.IsDev {
		setupVisualizer(mux, cfg.App.Host, cfg.App.Port)
	}
	if err := server.Run(context.Background(), mux, s); err != nil {
		log.Fatal(err)
	}
}
func main() {
	loadOptions := config.LoadOptions{
		RequireDotEnv: false,
	}
	config, err := config.Load(loadOptions)
	if err != nil {
		panic(`App Failed to load config ` + err.Error())
	}
	runServer(config)
}
