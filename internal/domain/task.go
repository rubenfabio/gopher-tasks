package domain

import "time"

// Task representa uma tarefa do usuário.
type Task struct {
    ID          string    // UUID gerado
    Title       string
    Description string
    DueDate     time.Time
    Completed   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
