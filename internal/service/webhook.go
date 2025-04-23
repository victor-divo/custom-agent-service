package service

import (
	"context"
	"log/slog"

	"github.com/victor-divo/custom-agent-service/internal/model"
	"github.com/victor-divo/custom-agent-service/internal/repository"
)

type WebhookService struct {
	Queue  repository.Queue
	Logger *slog.Logger
}

func NewWebhookService(queue repository.Queue, logger *slog.Logger) *WebhookService {
	return &WebhookService{
		Queue:  queue,
		Logger: logger,
	}
}

func (s *WebhookService) HandleWebhook(ctx context.Context, payload model.WebhookPayload) error {
	err := s.Queue.Enqueue(ctx, payload)
	if err != nil {
		s.Logger.Info("Webhook queued")
	}
	return err
}
