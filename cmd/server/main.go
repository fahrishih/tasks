package main

import (
	"log"

	"github.com/fahrishih/tasks/internal/task/adapters/inbound/ginhttp"
	"github.com/fahrishih/tasks/internal/task/adapters/outbound/memory"
	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/gin-gonic/gin"
)

func main() {
	// outbound adapters
	repo := memory.New()

	// application
	svc := app.NewService(repo)

	// inbound adapters
	handler := ginhttp.NewHandler(svc)

	r := gin.Default()
	handler.Register(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
