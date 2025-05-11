package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/config"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/database"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
)

func main() {
    // 1. Config
    cfg, err := config.Load("configs/config.yaml")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
        os.Exit(1)
    }

    // 2. Logger
    log := logger.New(cfg.Log.Level, cfg.Log.Format, os.Stdout)
    log.Info("Configuration loaded")

    // 3. DB
    db, err := database.Open(
        cfg.Database.Driver,
        cfg.Database.DSN,
        10,             // MaxOpenConns
        5,              // MaxIdleConns
        time.Minute*5,  // ConnMaxLifetime
    )
    if err != nil {
        log.WithField("error", err).Fatal("Failed to connect to database")
    }
    log.Info("Database connection established")

    // 4. Health check handler usando db.Ping()
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            log.WithField("error", err).Error("Database ping failed")
            http.Error(w, "service unavailable", http.StatusServiceUnavailable)
            return
        }
        log.WithField("path", r.URL.Path).Info("health check OK")
        w.Write([]byte("gopher-tasks is running and DB is healthy!"))
    })

    // 5. HTTP Server
    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    srv := &http.Server{
        Addr:         addr,
        Handler:      mux,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    log.Infof("Starting server on %s", addr)
    if err := srv.ListenAndServe(); err != nil {
        log.WithField("error", err).Fatal("Server failed")
    }
}
