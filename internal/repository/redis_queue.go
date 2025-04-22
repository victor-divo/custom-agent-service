package repository

import (
	"context"
	"encoding/json"
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

	return r.client.LPush(ctx, r.QueueName, payloadJSON).Err()
}

func (q *RedisQueue) Dequeue(ctx context.Context) (*model.WebhookPayload, error) {
	result, err := q.client.BRPop(ctx, 5*time.Second, q.QueueName).Result()
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
	if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}
