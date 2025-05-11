package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/config"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
)

func main() {
    // 1. Carrega config
    cfg, err := config.Load("configs/config.yaml")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
        os.Exit(1)
    }

    // 2. Inicializa logger com config
    log := logger.New(cfg.Log.Level, cfg.Log.Format, os.Stdout)
    log.Info("Configuration loaded")

    // 3. Health check handler
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.WithField("path", r.URL.Path).Info("health check")
        w.Write([]byte("gopher-tasks is running!"))
    })

    // 4. Server com timeouts
    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    srv := &http.Server{
        Addr:         addr,
        Handler:      mux,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    // use Infof para formatação e WithField para erro
    log.Infof("Starting server on %s", addr)
    if err := srv.ListenAndServe(); err != nil {
        log.WithField("error", err).Fatal("Server failed")
    }
}
