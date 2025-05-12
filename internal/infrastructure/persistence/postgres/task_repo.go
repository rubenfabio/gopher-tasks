package postgres

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rubenfabio/gopher-tasks/internal/domain"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(t *domain.Task) error {
	query := `
      INSERT INTO tasks (id, title, description, due_date, completed, created_at, updated_at)
      VALUES ($1,$2,$3,$4,$5,$6,$7)
    `
	now := time.Now()
	t.ID = uuid.NewString() // gera um UUID v4
	t.CreatedAt = now
	t.UpdatedAt = now

	_, err := r.db.Exec(query,
		t.ID,
		t.Title,
		t.Description,
		t.DueDate,
		t.Completed,
		t.CreatedAt,
		t.UpdatedAt,
	)
	return err
}
