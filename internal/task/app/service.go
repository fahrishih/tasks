package app

import (
	"context"
	"sync"

	"github.com/fahrishih/tasks/internal/task/domain"
	"github.com/google/uuid"
)

type Service struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*domain.Task
}

func NewService() *Service {
	return &Service{
		tasks: make(map[uuid.UUID]*domain.Task),
	}
}

// CreateTaskInput is the use case's input shape.
// It belongs to the app layer - not the domain (which doesn't know about
// "create requests") and not the transport (which has its own DTOs).
type CreateTaskInput struct {
	Title       string
	Description string
}

func (s *Service) CreateTask(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
	t, err := domain.NewTask(in.Title, in.Description)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.tasks[t.ID] = t
	s.mu.Unlock()
	return t, nil
}

func (s *Service) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	return t, nil
}

func (s *Service) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*domain.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	return out, nil
}

func (s *Service) CompleteTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	t.MarkComplete()
	return t, nil
}

func (s *Service) DeleteTask(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(s.tasks, id)
	return nil
}
