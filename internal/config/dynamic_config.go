package config

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type DynamicConfig struct {
	redis          *redis.Client
	logger         *slog.Logger
	maxAgentChat   int
	defaultMaxChat int
	mu             sync.RWMutex
	interval       time.Duration
}

func NewDynamicConfig(redis *redis.Client, logger *slog.Logger, interval time.Duration, fallback int) *DynamicConfig {
	dc := &DynamicConfig{
		redis:          redis,
		logger:         logger,
		interval:       interval,
		defaultMaxChat: fallback,
	}
	go dc.autoReload()
	return dc
}

func (c *DynamicConfig) autoReload() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		ctx := context.Background()
		val, err := c.redis.Get(ctx, "config:max_agent_chat").Result()
		if err != nil {
			c.logger.Error("Error reload config:max_agent_chat from Redis", "error", err)
			continue
		}

		maxAgentChat, err := strconv.Atoi(val)
		if err != nil {
			c.logger.Error("Invalid value for config:max_agent_chat", "value", val)
			continue
		}

		c.mu.Lock()
		c.maxAgentChat = maxAgentChat
		c.mu.Unlock()

		c.logger.Info("Dynamic config reloaded from redis", "config:max_agent_chat", maxAgentChat)
	}
}

func (c *DynamicConfig) GetMaxAgentChat() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.maxAgentChat == 0 {
		return c.defaultMaxChat
	}

	return c.maxAgentChat
}
