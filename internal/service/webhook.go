package service

import (
	"context"
	"log/slog"

	"github.com/victor-divo/custom-agent-service/internal/model"
	"github.com/victor-divo/custom-agent-service/internal/repository"
)

type WebhookService struct {
	Queue           repository.Queue
	Logger          *slog.Logger
	AgentRepository repository.AgentRepository
}

func NewWebhookService(queue repository.Queue, logger *slog.Logger, agentRepository repository.AgentRepository) *WebhookService {
	return &WebhookService{
		Queue:           queue,
		Logger:          logger,
		AgentRepository: agentRepository,
	}
}

func (s *WebhookService) HandleWebhook(ctx context.Context, payload model.WebhookPayload) error {
	err := s.Queue.Enqueue(ctx, payload)
	if err != nil {
		s.Logger.Info("Webhook queued")
	}
	return err
}

func (s *WebhookService) HandleResolvedChat(ctx context.Context, payload model.ResolvePayload) error {
	return s.AgentRepository.DecreaseCustomerCount(ctx, payload.ResolvedBy.ID)
}
