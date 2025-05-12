package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rubenfabio/gopher-tasks/internal/domain"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
	"github.com/rubenfabio/gopher-tasks/internal/usecase"
)

// createTaskRequest representa o payload para criação de Task.
type createTaskRequest struct {
    Title       string `json:"title" example:"Testar API"`
    Description string `json:"description" example:"Descrição da tarefa"`
    DueDate     string `json:"due_date" example:"2025-05-11T12:00:00Z"`
}

// TaskHandler agrupa os use cases e o logger para endpoints de Task.
type TaskHandler struct {
    CreateUC *usecase.CreateTaskUseCase
    ListUC   *usecase.ListTasksUseCase
    Log      logger.Logger
}

// NewTaskHandler injeta os use cases de criação e listagem, além do logger.
func NewTaskHandler(
    createUC *usecase.CreateTaskUseCase,
    listUC *usecase.ListTasksUseCase,
    log logger.Logger,
) *TaskHandler {
    return &TaskHandler{CreateUC: createUC, ListUC: listUC, Log: log}
}

// CreateTask godoc
// @Summary      Cria uma nova task
// @Description  Cria uma task com título, descrição e data de vencimento
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      createTaskRequest  true  "Payload para criar task"
// @Success      201   {object}  domain.Task
// @Failure      400   {object}  string
// @Failure      500   {object}  string
// @Router       /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req createTaskRequest
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

// ListTasks godoc
// @Summary      Lista tasks
// @Description  Retorna lista de tasks com filtros opcionais
// @Tags         tasks
// @Produce      json
// @Param        completed  query     bool   false  "Filtrar por concluídas"
// @Param        limit      query     int    false  "Limite de resultados"
// @Param        offset     query     int    false  "Offset para paginação"
// @Success      200        {array}   domain.Task
// @Failure      500        {object}  string
// @Router       /tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    var completed *bool
    if v := q.Get("completed"); v != "" {
        b, err := strconv.ParseBool(v)
        if err != nil {
            http.Error(w, "invalid completed filter", http.StatusBadRequest)
            return
        }
        completed = &b
    }

    limit := 0
    if v := q.Get("limit"); v != "" {
        if l, err := strconv.Atoi(v); err == nil {
            limit = l
        }
    }

    offset := 0
    if v := q.Get("offset"); v != "" {
        if o, err := strconv.Atoi(v); err == nil {
            offset = o
        }
    }

    filter := domain.TaskFilter{
        Completed: completed,
        Limit:     limit,
        Offset:    offset,
    }

    tasks, err := h.ListUC.Execute(filter)
    if err != nil {
        h.Log.WithField("error", err).Error("failed to list tasks")
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    // Garante que nunca seja retornado null, apenas um array vazio
    if tasks == nil {
        tasks = make([]*domain.Task, 0)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}
