package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/fahrishih/tasks/internal/task/domain"
)

// Repository is a Postgres implementation of app.Repository.
type Repository struct {
	pool *pgxpool.Pool
}

// New opens a Postgres connection pool and verifies it with a ping.
// The caller is responsible for calling Close() on the returned repo.
func New(ctx context.Context, dsn string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}
	return &Repository{pool: pool}, nil
}

// Close releases all pool connections. Call once at shutdown.
func (r *Repository) Close() {
	r.pool.Close()
}

// Compile-time assertion: the Postgres repo still satisfies app.Repository.
var _ app.Repository = (*Repository)(nil)

// Postgres SQLSTATE for unique_violation. We map it to ErrTaskDuplicate
// so the application doesn't need to know about Postgres error codes.
const sqlstateUniqueViolation = "23505"

func (r *Repository) Create(ctx context.Context, t *domain.Task) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tasks (id, title, description, completed, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, t.ID, t.Title, t.Description, t.Completed, t.CreatedAt)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == sqlstateUniqueViolation {
		return domain.ErrTaskDuplicate
	}
	if err != nil {
		return fmt.Errorf("postgres: create task: %w", err)
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, title, description, completed, created_at
		FROM tasks
		WHERE id = $1
	`, id)

	t := &domain.Task{}
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, domain.ErrTaskNotFound
	case err != nil:
		return nil, fmt.Errorf("postgres: get task: %w", err)
	}
	return t, nil
}

func (r *Repository) List(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, completed, created_at
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("postgres: list tasks: %w", err)
	}
	defer rows.Close()

	var out []*domain.Task
	for rows.Next() {
		t := &domain.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres: scan task: %w", err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: list rows: %w", err)
	}
	return out, nil
}

func (r *Repository) Update(ctx context.Context, t *domain.Task) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE tasks
		SET title = $2, description = $3, completed = $4
		WHERE id = $1
	`, t.ID, t.Title, t.Description, t.Completed)
	if err != nil {
		return fmt.Errorf("postgres: update task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("postgres: delete task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}
