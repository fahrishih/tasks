package ginhttp

import (
	"time"

	"github.com/fahrishih/tasks/internal/task/domain"
)

// createTaskRequest is the JSON body for POST /tasks.
// Gin's `binding` tags drive c.ShouldBindJSON validation.
type createTaskRequest struct {
	Title       string `json:"title"       binding:"required,max=200"`
	Description string `json:"description" binding:"max=2000"`
}

// taskResponse is the JSON shape returned to clients.
type taskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
}

func toResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
	}
}
