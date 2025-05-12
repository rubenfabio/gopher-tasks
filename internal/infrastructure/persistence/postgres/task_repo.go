package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

// Create insere uma nova Task no banco.
func (r *TaskRepo) Create(t *domain.Task) error {
    query := `
        INSERT INTO tasks (id, title, description, due_date, completed, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    now := time.Now()
    t.ID = uuid.NewString()
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

// FindByID busca uma Task pelo ID.
func (r *TaskRepo) FindByID(id string) (*domain.Task, error) {
    query := `
        SELECT id, title, description, due_date, completed, created_at, updated_at
        FROM tasks WHERE id = $1
    `
    row := r.db.QueryRow(query, id)
    var t domain.Task
    if err := row.Scan(
        &t.ID,
        &t.Title,
        &t.Description,
        &t.DueDate,
        &t.Completed,
        &t.CreatedAt,
        &t.UpdatedAt,
    ); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return &t, nil
}

// Update altera os campos de uma Task existente.
func (r *TaskRepo) Update(t *domain.Task) error {
    query := `
        UPDATE tasks
        SET title = $1, description = $2, due_date = $3, completed = $4, updated_at = $5
        WHERE id = $6
    `
    t.UpdatedAt = time.Now()
    _, err := r.db.Exec(query,
        t.Title,
        t.Description,
        t.DueDate,
        t.Completed,
        t.UpdatedAt,
        t.ID,
    )
    return err
}

// Delete remove uma Task pelo ID.
func (r *TaskRepo) Delete(id string) error {
    query := `DELETE FROM tasks WHERE id = $1`  
    res, err := r.db.Exec(query, id)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return errors.New("task not found")
    }
    return nil
}

// List retorna uma lista de Tasks segundo o filtro.
func (r *TaskRepo) List(filter domain.TaskFilter) ([]*domain.Task, error) {
    query := `
        SELECT id, title, description, due_date, completed, created_at, updated_at
        FROM tasks
    `
    var args []interface{}
    var conditions []string

    if filter.Completed != nil {
        conditions = append(conditions, "completed = $1")
        args = append(args, *filter.Completed)
    }
    if len(conditions) > 0 {
        query += " WHERE " + strings.Join(conditions, " AND ")
    }
    query += " ORDER BY created_at DESC"
    if filter.Limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", filter.Limit)
    }
    if filter.Offset > 0 {
        query += fmt.Sprintf(" OFFSET %d", filter.Offset)
    }

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*domain.Task
    for rows.Next() {
        var t domain.Task
        if err := rows.Scan(
            &t.ID,
            &t.Title,
            &t.Description,
            &t.DueDate,
            &t.Completed,
            &t.CreatedAt,
            &t.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        tasks = append(tasks, &t)
    }
    return tasks, nil
}
