package worker

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/victor-divo/custom-agent-service/internal/config"
	"github.com/victor-divo/custom-agent-service/internal/model"
	"github.com/victor-divo/custom-agent-service/internal/repository"
	"github.com/victor-divo/custom-agent-service/internal/service"
)

type WebhookWorker struct {
	Queue           repository.Queue
	Logger          *slog.Logger
	DynamicConfig   *config.DynamicConfig
	AgentRepository *repository.AgentRepository
}

func NewWebhookWorker(queue repository.Queue, logger *slog.Logger, DynamicConfig *config.DynamicConfig, AgentRepository *repository.AgentRepository) *WebhookWorker {
	return &WebhookWorker{
		Queue:           queue,
		Logger:          logger,
		DynamicConfig:   DynamicConfig,
		AgentRepository: AgentRepository,
	}
}

func (w *WebhookWorker) Start(ctx context.Context) {

	go func() {
		for {
			payload, err := w.Queue.Dequeue(ctx)
			if err != nil {
				w.Logger.Error("Error dequeuing payload:", "error", err)
				continue
			}

			if payload == nil {
				w.Logger.Info("No payload to process")
				time.Sleep(2 * time.Second)
				continue
			}

			isAgentFull, err := w.processPayload(payload)
			if err != nil {
				w.Logger.Error("Error processing payload:", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}
			if isAgentFull {
				w.Logger.Info("All agents are busy")
				err = w.Queue.Requeue(ctx, *payload)
				if err != nil {
					w.Logger.Error("Error requeuing payload:", "error", err)
				}
				time.Sleep(5 * time.Second)
				continue
			}
			w.Logger.Info("Message allocated successfully")
		}
	}()
}

func (w *WebhookWorker) processPayload(payload *model.WebhookPayload) (bool, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("processing data")
	agentManagementService := service.NewAgentManagementService(logger)
	var agents []model.Agent

	localDataAgents, err := w.AgentRepository.GetAllAgents(context.Background())
	if err != nil {
		logger.Error("Error fetching local agents:", "error", err)
	}

	if len(localDataAgents) != 0 {
		agents = localDataAgents
	} else {
		agents, err = agentManagementService.GetAllAgents(context.Background())
		if err != nil {
			logger.Error("Error fetching agents:", "error", err)
			return false, err
		}
		w.AgentRepository.SetInitialAgent(context.Background(), agents)
	}

	isAgentFull := false
	for i, agent := range agents {
		if isAgentEnglible(agent, w) {
			agentManagementService.AssignAgent(context.Background(), agent.ID, payload.RoomID)
			err = w.AgentRepository.IncreaseCustomerCount(context.Background(), agent.ID)
			if err != nil {
				logger.Error("Error increasing customer count:", "error", err)
				return false, err
			}
			break
		}
		if i == len(agents)-1 {
			isAgentFull = true
		}
	}

	return isAgentFull, nil
}

func isAgentEnglible(agent model.Agent, w *WebhookWorker) bool {
	return agent.CurrentCustomerCount < w.DynamicConfig.GetMaxAgentChat()
}
