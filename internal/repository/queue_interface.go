package repository

import (
	"context"

	"github.com/victor-divo/custom-agent-service/internal/model"
)

type Queue interface {
	Enqueue(ctx context.Context, payload model.WebhookPayload) error
	Dequeue(ctx context.Context) (*model.WebhookPayload, error)
	Requeue(ctx context.Context, payload model.WebhookPayload) error
}
