package repository

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(redisAddr string) *RedisRepository {
	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	err := client.XGroupCreateMkStream(context.Background(), STREAM_REQUEST, CONSUMER_GROUP, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	return &RedisRepository{Client: client}
}

func (r *RedisRepository) ConsumeMessages(ctx context.Context, consumerName string) (*redis.XMessage, error) {
	streams, err := r.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    CONSUMER_GROUP,
		Consumer: consumerName,
		Streams:  []string{STREAM_REQUEST, ">"},
		Count:    1,
		Block:    0,
	}).Result()

	if err != nil || len(streams) == 0 || len(streams[0].Messages) == 0 {
		return nil, err
	}
	return &streams[0].Messages[0], nil
}

func (r *RedisRepository) PublishSearchResult(ctx context.Context, searchID, resultsJSON string) error {
	message := map[string]interface{}{
		"search_id": searchID,
		"status":    "completed",
		"results":   resultsJSON,
	}
	return r.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: STREAM_RESULT,
		Values: message,
	}).Err()
}

func (r *RedisRepository) AcknowledgeMessage(ctx context.Context, messageID string) error {
	return r.Client.XAck(ctx, STREAM_REQUEST, CONSUMER_GROUP, messageID).Err()
}
