package ginhttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/fahrishih/tasks/internal/task/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler is the inbound HTTP adapter. It depends on app.Service
// (the application boundary) and knows nothing about storage.
type Handler struct {
	svc *app.Service
}

func NewHandler(svc *app.Service) *Handler {
	return &Handler{svc: svc}
}

// Register attaches all task routes onto the given Gin engine.
// The engine itself is owned by the composition root.
func (h *Handler) Register(r *gin.Engine) {
	r.POST("/tasks", h.createTask)
	r.GET("/tasks/:id", h.getTask)
	r.GET("/tasks", h.listTasks)
	r.PATCH("/tasks/:id", h.completeTask)
	r.DELETE("/tasks/:id", h.deleteTask)
}

// ---------- DTOs ----------
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

// ---------- error mapping ----------
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

// ---------- handlers ----------

func (h *Handler) createTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.svc.CreateTask(c.Request.Context(), app.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toResponse(t))
}

func (h *Handler) getTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	t, err := h.svc.GetTask(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(t))
}

func (h *Handler) listTasks(c *gin.Context) {
	ts, err := h.svc.ListTasks(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	out := make([]taskResponse, 0, len(ts))
	for _, t := range ts {
		out = append(out, toResponse(t))
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) completeTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	t, err := h.svc.CompleteTask(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(t))
}

func (h *Handler) deleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeleteTask(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
