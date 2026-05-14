package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	ID          string    `json:id`
	Title       string    `json:title`
	Description string    `json:description`
	Completed   bool      `json:completed`
	CreatedAt   time.Time `string:createdAt`
}

var (
	tasks = make(map[string]Task)
	mu    sync.Mutex
)

func main() {
	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		var t Task
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if t.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
			return
		}
		t.ID = uuid.NewString()
		t.CreatedAt = time.Now()
		t.Completed = false

		mu.Lock()
		tasks[t.ID] = t
		mu.Unlock()

		c.JSON(http.StatusCreated, t)
	})

	r.GET("/tasks/:id", func(c *gin.Context) {
		mu.Lock()
		t, ok := tasks[c.Param("id")]
		mu.Unlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	r.GET("/tasks", func(c *gin.Context) {
		mu.Lock()
		out := make([]Task, 0, len(tasks))
		for _, t := range tasks {
			out = append(out, t)
		}
		mu.Unlock()
		c.JSON(http.StatusOK, out)
	})

	r.PATCH("/tasks/:id", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		t, ok := tasks[c.Param("id")]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "tasks not found"})
			return
		}
		t.Completed = true
		tasks[t.ID] = t
		c.JSON(http.StatusOK, t)
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		mu.Lock()
		_, ok := tasks[c.Param("id")]
		if ok {
			delete(tasks, c.Param("id"))
		}
		mu.Unlock()
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.Run(":8080")
}
