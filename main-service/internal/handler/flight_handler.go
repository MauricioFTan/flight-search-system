package handler

import (
	"bufio"
	"fmt"
	"log"

	"github.com/MauricioFTan/main-service/internal/model"
	"github.com/MauricioFTan/main-service/internal/repository"
	"github.com/MauricioFTan/main-service/internal/sse"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FlightHandler struct {
	Repo *repository.RedisRepository
	Hub  *sse.Hub
}

func NewFlightHandler(repo *repository.RedisRepository, hub *sse.Hub) *FlightHandler {
	return &FlightHandler{Repo: repo, Hub: hub}
}

func (h *FlightHandler) CreateSearch(c *fiber.Ctx) error {
	req := new(model.SearchRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	searchID := uuid.New().String()

	searchData := model.SearchRequestData{
		SearchID:   searchID,
		From:       req.From,
		To:         req.To,
		Date:       req.Date,
		Passengers: req.Passengers,
	}

	if err := h.Repo.PublishSearchRequest(c.Context(), searchData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to submit search request"})
	}

	response := model.SearchResponse{
		Data: model.SearchData{
			SearchID: searchID,
			Status:   "processing",
		},
		Success: true,
		Message: "Search request submitted",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *FlightHandler) StreamSearch(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	searchID := c.Params("search_id")

	clientClosed := c.Context().Done()

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {

		messageChan := h.Hub.Register(searchID)
		defer h.Hub.Unregister(searchID)

		for {
			select {
			case message, ok := <-messageChan:
				if !ok {
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", message)

				if err := w.Flush(); err != nil {
					log.Printf("Error flushing data for search_id %s: %v", searchID, err)
					return
				}

			case <-clientClosed:
				log.Printf("Client for search_id %s closed connection.", searchID)
				return
			}
		}
	})

	return nil
}
