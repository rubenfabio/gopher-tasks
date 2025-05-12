// @title       Gopher Tasks API
// @version     1.0
// @description API para gerenciamento de tarefas
// @host        localhost:8080
// @BasePath    /
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/rubenfabio/gopher-tasks/docs" // swagger docs
	httpdelivery "github.com/rubenfabio/gopher-tasks/internal/delivery/http"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/config"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/database"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/persistence/postgres"
	"github.com/rubenfabio/gopher-tasks/internal/usecase"
	httpSwagger "github.com/swaggo/http-swagger" // swagger UI handler
)

func main() {
    // 1. Carrega configuração
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
        10,            // MaxOpenConns
        5,             // MaxIdleConns
        time.Minute*5, // ConnMaxLifetime
    )
    if err != nil {
        log.WithField("error", err).Fatal("Failed to connect to database")
    }
    log.Info("Database connection established")

    // 4. UseCases e Handler
    taskRepo    := postgres.NewTaskRepo(db)
    createUC    := usecase.NewCreateTaskUseCase(taskRepo)
    listUC      := usecase.NewListTasksUseCase(taskRepo)
    taskHandler := httpdelivery.NewTaskHandler(createUC, listUC, log)

    // 5. Router
    r := mux.NewRouter()

    // Swagger UI endpoint em /swagger/index.html
    r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    // expõe swagger.json em /swagger.json
    r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "docs/swagger.json")
    }).Methods(http.MethodGet)
    // Health-check
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := db.Ping(); err != nil {
            log.WithField("error", err).Error("Database ping failed")
            http.Error(w, "service unavailable", http.StatusServiceUnavailable)
            return
        }
        log.WithField("path", r.URL.Path).Info("health check OK")
        w.Write([]byte("gopher-tasks is running and DB is healthy!"))
    }).Methods(http.MethodGet)

    // Create task
    r.HandleFunc("/tasks", taskHandler.Create).Methods(http.MethodPost)
    // List tasks
    r.HandleFunc("/tasks", taskHandler.List).Methods(http.MethodGet)

    // 6. Start server
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
