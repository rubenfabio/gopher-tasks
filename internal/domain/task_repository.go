package domain

// TaskRepository define as operações de persistência de Task.
type TaskRepository interface {
    Create(task *Task) error
    FindByID(id string) (*Task, error)
    Update(task *Task) error
    Delete(id string) error
    List(filter TaskFilter) ([]*Task, error)
}

// TaskFilter para paginação/filtros
type TaskFilter struct {
    Completed *bool
    Limit     int
    Offset    int
}
