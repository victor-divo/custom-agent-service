package repository

import (
	"context"
	"sort"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/victor-divo/custom-agent-service/internal/model"
)

type AgentRepository struct {
	client  *redis.Client
	KeyName string
}

func NewAgentRepository(client *redis.Client, keyName string) *AgentRepository {
	return &AgentRepository{
		client:  client,
		KeyName: keyName,
	}
}

func (r *AgentRepository) SetInitialAgent(ctx context.Context, agents []model.Agent) error {

	args := make([]interface{}, 0, len(agents)*2)
	for _, agent := range agents {
		args = append(args, agent.ID, agent.CurrentCustomerCount)
	}

	return r.client.HSet(ctx, r.KeyName, args...).Err()
}

func (r *AgentRepository) IncreaseCustomerCount(ctx context.Context, agentID int) error {
	stringID := strconv.Itoa(agentID)
	_, err := r.client.HIncrBy(ctx, r.KeyName, stringID, 1).Result()
	return err
}

func (r *AgentRepository) DecreaseCustomerCount(ctx context.Context, agentID int) error {
	stringID := strconv.Itoa(agentID)
	_, err := r.client.HIncrBy(ctx, r.KeyName, stringID, -1).Result()
	return err
}

func (r *AgentRepository) GetAllAgents(ctx context.Context) ([]model.Agent, error) {
	result, err := r.client.HGetAll(ctx, r.KeyName).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	agents := make([]model.Agent, 0, len(result))
	for key, value := range result {
		// Convert the value from string to int
		intCustCount, _ := strconv.Atoi(value)
		intID, _ := strconv.Atoi(key)

		agent := model.Agent{
			ID:                   intID,
			CurrentCustomerCount: intCustCount,
		}
		agents = append(agents, agent)
	}

	sort.Slice(agents, func(i, j int) bool {
		return agents[i].CurrentCustomerCount < agents[j].CurrentCustomerCount
	})

	return agents, nil
}
