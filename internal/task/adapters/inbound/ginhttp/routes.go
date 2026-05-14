package ginhttp

import (
	"net/http"

	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
