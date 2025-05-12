package usecase

import "github.com/rubenfabio/gopher-tasks/internal/domain"

// ListTasksUseCase encapsula a lógica de listar Tasks.
type ListTasksUseCase struct {
    Repo domain.TaskRepository
}

func NewListTasksUseCase(repo domain.TaskRepository) *ListTasksUseCase {
    return &ListTasksUseCase{Repo: repo}
}

// Execute retorna as tasks de acordo com o filtro.
func (uc *ListTasksUseCase) Execute(filter domain.TaskFilter) ([]*domain.Task, error) {
    return uc.Repo.List(filter)
}
