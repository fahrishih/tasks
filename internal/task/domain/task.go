package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
}

// NewTask is a constructor that enforces invariants.
// A Task can only enter the system through this function
func NewTask(title, description string) (*Task, error) {
	t := &Task{
		ID:          uuid.New(),
		Title:       strings.TrimSpace(title),
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
	}
	if err := t.Validate(); err != nil {
		return nil, err
	}
	return t, nil
}

// Validate check the domain invariants
func (t *Task) Validate() error {
	if t.Title == "" {
		return ErrInvalidTitle
	}
	if len(t.Title) > 200 {
		return ErrTitleTooLong
	}
	return nil
}

// MarkComplete is a domain behaviour - not just a setter.
// Putting this here means the rule "completeing is a one-way operation"
// can be enforced in one place
func (t *Task) MarkComplete() {
	t.Completed = true
}
