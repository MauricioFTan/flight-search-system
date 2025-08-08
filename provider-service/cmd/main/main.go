package main

import (
	"context"
	"os"

	"github.com/MauricioFTan/provider-service/internal/repository"
	"github.com/MauricioFTan/provider-service/internal/worker"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	repo := repository.NewRedisRepository(redisAddr)

	consumer := worker.NewConsumer(repo, "processor-1")

	consumer.Start(context.Background())
}
