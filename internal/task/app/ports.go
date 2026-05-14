package app

import (
	"context"

	"github.com/fahrishih/tasks/internal/task/domain"
	"github.com/google/uuid"
)

// Repository is the persistence port the application requires.
// the application DEFINES it; outbound adapters IMPLEMENT it.
//
// Contract:
// - Get / Update / Delete must return domain.ErrTaskNotFound when no task matches the given id.
// - Create must return domain.ErrTaskDuplicate if a task with the same id already exists
type Repository interface {
	Create(ctx context.Context, t *domain.Task) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	List(ctx context.Context) ([]*domain.Task, error)
	Update(ctx context.Context, t *domain.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
