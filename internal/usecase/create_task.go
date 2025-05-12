package usecase

import (
	"time"

	"github.com/rubenfabio/gopher-tasks/internal/domain"
)

// CreateTaskUseCase encapsula a lógica de criar uma Task.
type CreateTaskUseCase struct {
    Repo domain.TaskRepository
}

// NewCreateTaskUseCase injeta o repositório de tarefas.
func NewCreateTaskUseCase(repo domain.TaskRepository) *CreateTaskUseCase {
    return &CreateTaskUseCase{Repo: repo}
}

// Execute cria uma nova Task no repositório e retorna a entidade preenchida.
func (uc *CreateTaskUseCase) Execute(title, description string, dueDate time.Time) (*domain.Task, error) {
    task := &domain.Task{
        Title:       title,
        Description: description,
        DueDate:     dueDate,
    }
    if err := uc.Repo.Create(task); err != nil {
        return nil, err
    }
    return task, nil
}
