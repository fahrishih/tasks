package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/fahrishih/tasks/internal/task/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// -- Request / Response shapes (will move to a DTO file in stage 5)

type createTaskRequest struct {
	Title       string `json:title`
	Description string `json:description`
}

type taskResponse struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Desription string    `json:"description"`
	Completed  bool      `json:"completed"`
	CreatedAt  time.Time `json:"createdAt"`
}

func toResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:         t.ID.String(),
		Title:      t.Title,
		Desription: t.Description,
		Completed:  t.Completed,
		CreatedAt:  t.CreatedAt,
	}
}

var (
	tasks = make(map[uuid.UUID]*domain.Task)
	mu    sync.Mutex
)

func main() {
	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		var req createTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		t, err := domain.NewTask(req.Title, req.Description)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mu.Lock()
		tasks[t.ID] = t
		mu.Unlock()

		c.JSON(http.StatusCreated, toResponse(t))
	})

	r.GET("/tasks/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		mu.Lock()
		t, ok := tasks[id]
		mu.Unlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrTaskNotFound})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	r.GET("/tasks", func(c *gin.Context) {
		mu.Lock()
		out := make([]taskResponse, 0, len(tasks))
		for _, t := range tasks {
			out = append(out, toResponse(t))
		}
		mu.Unlock()
		c.JSON(http.StatusOK, out)
	})

	r.PATCH("/tasks/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		mu.Lock()
		defer mu.Unlock()
		t, ok := tasks[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrTaskNotFound})
			return
		}
		t.MarkComplete()
		c.JSON(http.StatusOK, toResponse(t))
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		mu.Lock()
		_, ok := tasks[id]
		if ok {
			delete(tasks, id)
		}
		mu.Unlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrTaskNotFound})
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.Run(":8080")
}
