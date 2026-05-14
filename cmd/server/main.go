package main

import (
	"context"
	"log"
	"os"

	"github.com/fahrishih/tasks/internal/task/adapters/inbound/ginhttp"
	"github.com/fahrishih/tasks/internal/task/adapters/outbound/postgres"
	"github.com/fahrishih/tasks/internal/task/app"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// outbound adapters
	repo, err := postgres.New(ctx, dsn)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer repo.Close()

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
