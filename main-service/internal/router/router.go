package router

import (
	"github.com/MauricioFTan/main-service/internal/handler"
	"github.com/MauricioFTan/main-service/internal/repository"
	"github.com/MauricioFTan/main-service/internal/sse"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, redisRepo *repository.RedisRepository, hub *sse.Hub) {
	app.Use(logger.New())

	flightHandler := handler.NewFlightHandler(redisRepo, hub)

	api := app.Group("/api")
	flights := api.Group("/flights")

	flights.Post("/search", flightHandler.CreateSearch)

	flights.Get("/search/:search_id/stream", flightHandler.StreamSearch)
}
