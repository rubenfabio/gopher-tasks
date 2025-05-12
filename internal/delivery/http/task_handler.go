package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
	"github.com/rubenfabio/gopher-tasks/internal/usecase"
)

type TaskHandler struct {
    CreateUC *usecase.CreateTaskUseCase
    Log      logger.Logger
}

// NewTaskHandler injeta o use case e o logger.
func NewTaskHandler(uc *usecase.CreateTaskUseCase, log logger.Logger) *TaskHandler {
    return &TaskHandler{CreateUC: uc, Log: log}
}

// Create endpoint: POST /tasks
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        DueDate     string `json:"due_date"` // ISO8601: "2025-05-11T12:00:00Z"
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }

    due, err := time.Parse(time.RFC3339, req.DueDate)
    if err != nil {
        http.Error(w, "invalid due_date format", http.StatusBadRequest)
        return
    }

    task, err := h.CreateUC.Execute(req.Title, req.Description, due)
    if err != nil {
        h.Log.WithField("error", err).Error("failed to create task")
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}
