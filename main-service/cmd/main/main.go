package main

import (
	"context"
	"log"
	"os"

	"github.com/MauricioFTan/main-service/internal/repository"
	"github.com/MauricioFTan/main-service/internal/router"
	"github.com/MauricioFTan/main-service/internal/sse"

	"github.com/gofiber/fiber/v2"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisRepo := repository.NewRedisRepository(redisAddr)

	sseHub := sse.NewHub()

	go redisRepo.ListenForResults(context.Background(), sseHub)

	app := fiber.New()

	router.SetupRoutes(app, redisRepo, sseHub)

	log.Println("Main service is starting on port 8080...")
	app.Listen(":8080")
}
