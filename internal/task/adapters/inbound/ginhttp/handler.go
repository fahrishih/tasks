package ginhttp

import (
	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/gin-gonic/gin"
)

// Handler is the inbound HTTP adapter. It depends on app.Service
// (the application boundary) and knows nothing about storage.
type Handler struct {
	svc *app.Service
}

func NewHandler(svc *app.Service) *Handler {
	return &Handler{svc: svc}
}

// Register attaches all task routes under /api/v1 onto the given engine.
// The engine itself is owned by the composition root.
func (h *Handler) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/tasks", h.createTask)
		v1.GET("/tasks/:id", h.getTask)
		v1.GET("/tasks", h.listTasks)
		v1.PATCH("/tasks/:id", h.completeTask)
		v1.DELETE("/tasks/:id", h.deleteTask)
	}
}
