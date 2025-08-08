package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MauricioFTan/main-service/internal/model"
	"github.com/MauricioFTan/main-service/internal/sse"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(redisAddr string) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return &RedisRepository{Client: client}
}

func (r *RedisRepository) PublishSearchRequest(ctx context.Context, data model.SearchRequestData) error {

	message := map[string]interface{}{
		"search_id":  data.SearchID,
		"from":       data.From,
		"to":         data.To,
		"date":       data.Date,
		"passengers": data.Passengers,
	}

	err := r.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: STREAM_REQUEST,
		Values: message,
	}).Err()

	if err != nil {
		log.Printf("Failed to publish search request to Redis: %v", err)
		return err
	}

	log.Printf("Successfully published search ID %s to stream %s", data.SearchID, STREAM_REQUEST)
	return nil
}

func (r *RedisRepository) ListenForResults(ctx context.Context, hub *sse.Hub) {
	lastID := "0-0"

	log.Printf("Listener started for stream: %s.", STREAM_RESULT)

	for {
		streams, err := r.Client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{STREAM_RESULT, lastID},
			Count:   1,
			Block:   0,
		}).Result()

		if err != nil {
			log.Printf("Error reading from Redis Stream '%s': %v", STREAM_RESULT, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(streams) == 0 || len(streams[0].Messages) == 0 {
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				lastID = message.ID

				searchID, ok1 := message.Values["search_id"].(string)
				results, ok2 := message.Values["results"].(string)

				if ok1 && ok2 {
					jsonMessage := fmt.Sprintf(`{"search_id":"%s", "status":"completed", "results":%s}`, searchID, results)
					hub.Broadcast(searchID, jsonMessage)
				} else {
					log.Printf("Message %s has invalid format.", message.ID)
				}
			}
		}
	}
}
