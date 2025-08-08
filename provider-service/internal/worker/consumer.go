package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MauricioFTan/provider-service/internal/repository"
)

type Consumer struct {
	Repo         *repository.RedisRepository
	ConsumerName string
}

func NewConsumer(repo *repository.RedisRepository, name string) *Consumer {
	return &Consumer{Repo: repo, ConsumerName: name}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Printf("Consumer '%s' started. Listening for messages...", c.ConsumerName)

	for {
		message, err := c.Repo.ConsumeMessages(ctx, c.ConsumerName)
		if err != nil {
			log.Printf("Error consuming messages: %v. Retrying in 2 seconds...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if message == nil {
			continue
		}

		log.Printf("Processing message ID: %s", message.ID)

		searchID, ok := message.Values["search_id"].(string)
		if !ok {
			log.Printf("Invalid message format, missing search_id. Acknowledging to discard.")
			c.Repo.AcknowledgeMessage(ctx, message.ID)
			continue
		}

		time.Sleep(30 * time.Second)
		mockResults := fmt.Sprintf(`[{"flight_number":"GA%d","price":1500000},{"flight_number":"JT%d","price":1200000}]`, time.Now().UnixMilli()%100, time.Now().UnixMilli()%200)

		if err := c.Repo.PublishSearchResult(ctx, searchID, mockResults); err != nil {
			log.Printf("Failed to publish results for %s: %v", searchID, err)

			continue
		}
		log.Printf("Published results for search ID: %s", searchID)

		if err := c.Repo.AcknowledgeMessage(ctx, message.ID); err != nil {
			log.Printf("Failed to acknowledge message %s: %v", message.ID, err)
		} else {
			log.Printf("Acknowledged message ID: %s", message.ID)
		}
	}
}
