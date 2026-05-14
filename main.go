package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/fahrishih/tasks/internal/task/adapters/outbound/memory"
	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/fahrishih/tasks/internal/task/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// -- Request / Response shapes (will move to a DTO file in stage 5)

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

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

// --- Domain error → HTTP status mapping (will move in Stage 5) ---

func statusForError(err error) int {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrTaskDuplicate):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalidTitle),
		errors.Is(err, domain.ErrTitleTooLong):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeError(c *gin.Context, err error) {
	c.JSON(statusForError(err), gin.H{"error": err.Error()})
}

func main() {
	repo := memory.New()
	svc := app.NewService(repo)
	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		var req createTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := svc.CreateTask(c.Request.Context(), app.CreateTaskInput{
			Title:       req.Title,
			Description: req.Description,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusCreated, toResponse(t))
	})

	r.GET("/tasks/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		t, err := svc.GetTask(c.Request.Context(), id)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, t)
	})

	r.GET("/tasks", func(c *gin.Context) {
		ts, err := svc.ListTasks(c.Request.Context())
		if err != nil {
			writeError(c, err)
			return
		}
		out := make([]taskResponse, 0, len(ts))
		for _, t := range ts {
			out = append(out, toResponse(t))
		}
		c.JSON(http.StatusOK, out)
	})

	r.PATCH("/tasks/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		t, err := svc.CompleteTask(c.Request.Context(), id)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, toResponse(t))
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := svc.DeleteTask(c.Request.Context(), id); err != nil {
			writeError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.Run(":8080")
}
