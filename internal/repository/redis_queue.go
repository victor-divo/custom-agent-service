package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/victor-divo/custom-agent-service/internal/model"
)

type RedisQueue struct {
	client    *redis.Client
	QueueName string
}

func NewRedisQueue(client *redis.Client, queueName string) *RedisQueue {
	return &RedisQueue{
		client:    client,
		QueueName: queueName,
	}
}

func (r *RedisQueue) Enqueue(ctx context.Context, payload model.WebhookPayload) error {
	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	exists, err := r.client.SIsMember(ctx, r.QueueName+"_set", payload.RoomID).Result()
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("payload with RoomID %s already exists in the queue", payload.RoomID)
	}

	err = r.client.SAdd(ctx, r.QueueName+"_set", payload.RoomID).Err()
	if err != nil {
		return fmt.Errorf("failed to add RoomID to set: %w", err)
	}

	return r.client.RPush(ctx, r.QueueName, payloadJSON).Err()
}

func (r *RedisQueue) Dequeue(ctx context.Context) (*model.WebhookPayload, error) {
	result, err := r.client.BLPop(ctx, 5*time.Second, r.QueueName).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	if len(result) < 2 {
		return nil, nil
	}

	var payload model.WebhookPayload
	err = json.Unmarshal([]byte(result[1]), &payload)
	if err != nil {
		return nil, err
	}

	_, err = r.client.SRem(ctx, r.QueueName+"_set", payload.RoomID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to remove RoomID from set: %w", err)
	}

	return &payload, nil
}

func (r *RedisQueue) Requeue(ctx context.Context, payload model.WebhookPayload) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = r.client.SAdd(ctx, r.QueueName+"_set", payload.RoomID).Err()
	if err != nil {
		return fmt.Errorf("failed to add RoomID to set: %w", err)
	}

	return r.client.LPush(ctx, r.QueueName, payloadJSON).Err()
}
