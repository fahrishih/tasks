package memory

import (
	"context"
	"sync"

	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/fahrishih/tasks/internal/task/domain"
	"github.com/google/uuid"
)

type Repository struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*domain.Task
}

func New() *Repository {
	return &Repository{
		tasks: make(map[uuid.UUID]*domain.Task),
	}
}

// Compile-time assertion that *Repository satisfies app.Repository.
// If you ever break the interface, the build will fail here.
var _ app.Repository = (*Repository)(nil)

func (r *Repository) Create(ctx context.Context, t *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[t.ID]; ok {
		return domain.ErrTaskDuplicate
	}
	r.tasks[t.ID] = t
	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	return t, nil
}

func (r *Repository) List(ctx context.Context) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		out = append(out, t)
	}
	return out, nil
}

func (r *Repository) Update(ctx context.Context, t *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[t.ID]; !ok {
		return domain.ErrTaskNotFound
	}
	r.tasks[t.ID] = t
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(r.tasks, id)
	return nil
}
