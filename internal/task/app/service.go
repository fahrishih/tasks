package app

import (
	"context"

	"github.com/fahrishih/tasks/internal/task/domain"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateTaskInput struct {
	Title       string
	Description string
}

func (s *Service) CreateTask(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
	t, err := domain.NewTask(in.Title, in.Description)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) ListTasks(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) CompleteTask(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	t.MarkComplete()
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
