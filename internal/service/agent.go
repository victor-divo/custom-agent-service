package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/victor-divo/custom-agent-service/internal/config"
	"github.com/victor-divo/custom-agent-service/internal/model"
)

type AgentManagementService struct {
	Client *resty.Client
	Logger *slog.Logger
}

func NewAgentManagementService(logger *slog.Logger) *AgentManagementService {
	client := resty.New().
		SetBaseURL(config.GetBaseURL()).
		SetHeader("Qiscus-Secret-Key", config.GetSecretKey()).
		SetHeader("Qiscus-App-Id", config.GetAppId()).
		SetHeader("Content-Type", "application/json")
	return &AgentManagementService{
		Client: client,
		Logger: logger,
	}
}

func (s *AgentManagementService) GetAllAgents(ctx context.Context) ([]model.Agent, error) {
	var getAllAgentsResponse model.GetAllDivisionAgentsResponse

	_, err := s.Client.R().
		SetContext(ctx).
		SetResult(&getAllAgentsResponse).
		Get("/api/v2/admin/agents/by_division?page=1&limit=100&division_ids[]=133584&is_available=true&sort=asc")

	if err != nil {
		s.Logger.Error("Error fetching agents:", "error", err)
		return nil, err
	}
	agents := getAllAgentsResponse.Data

	return agents, nil
}

func (s *AgentManagementService) AssignAgent(ctx context.Context, agentID int, roomID string) error {
	formData := map[string]string{
		"agent_id": strconv.Itoa(agentID),
		"room_id":  roomID,
	}

	var result struct {
		Status string                 `json:"status,omitempty"`
		Data   map[string]interface{} `json:"data,omitempty"`
		Error  string                 `json:"error,omitempty"`
	}

	resp, err := s.Client.R().
		SetContext(ctx).
		SetFormData(formData).
		SetResult(&result).
		Post("/api/v1/admin/service/assign_agent")

	if err != nil {
		s.Logger.Error("Error assigning agent to room (network error):", "error", err)
		return fmt.Errorf("network error: %w", err)
	}
	if !resp.IsSuccess() {
		s.Logger.Error("Failed to assign agent to room (API error)",
			"status_code", resp.StatusCode(),
			"response", resp.String(),
		)
		return fmt.Errorf("API error: status %d", resp.StatusCode())
	}
	s.Logger.Info("Agent successfully assigned to room", "room_id", roomID, "agent_id", agentID)
	return nil
}
