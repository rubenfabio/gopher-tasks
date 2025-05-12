package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	httpdelivery "github.com/rubenfabio/gopher-tasks/internal/delivery/http"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/config"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/database"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/persistence/postgres"
	"github.com/rubenfabio/gopher-tasks/internal/usecase"
)

func main() {
    // 1. Carrega configuração
    cfg, err := config.Load("configs/config.yaml")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
        os.Exit(1)
    }

    // 2. Inicializa logger
    log := logger.New(cfg.Log.Level, cfg.Log.Format, os.Stdout)
    log.Info("Configuration loaded")

    // 3. Abre conexão com DB
    db, err := database.Open(
        cfg.Database.Driver,
        cfg.Database.DSN,
        10,            // MaxOpenConns
        5,             // MaxIdleConns
        time.Minute*5, // ConnMaxLifetime
    )
    if err != nil {
        log.WithField("error", err).Fatal("Failed to connect to database")
    }
    log.Info("Database connection established")

    // 4. Inicializa repositório, caso de uso e handler
    taskRepo := postgres.NewTaskRepo(db)
    createUC := usecase.NewCreateTaskUseCase(taskRepo)
    taskHandler := httpdelivery.NewTaskHandler(createUC, log)

    // 5. Configura router
    r := mux.NewRouter()

    // Health-check endpoint
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            log.WithField("error", err).Error("Database ping failed")
            http.Error(w, "service unavailable", http.StatusServiceUnavailable)
            return
        }
        log.WithField("path", r.URL.Path).Info("health check OK")
        w.Write([]byte("gopher-tasks is running and DB is healthy!"))
    }).Methods(http.MethodGet)

    // Create task endpoint
    r.HandleFunc("/tasks", taskHandler.Create).Methods(http.MethodPost)

    // 6. Inicia servidor HTTP
    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    srv := &http.Server{
        Addr:         addr,
        Handler:      r,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }

    log.Infof("Starting server on %s", addr)
    if err := srv.ListenAndServe(); err != nil {
        log.WithField("error", err).Fatal("Server failed")
    }
}
