package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/victor-divo/custom-agent-service/internal/config"
	"github.com/victor-divo/custom-agent-service/internal/handler"
	"github.com/victor-divo/custom-agent-service/internal/repository"
	"github.com/victor-divo/custom-agent-service/internal/service"
	"github.com/victor-divo/custom-agent-service/internal/worker"
)

func main() {
	// set logger
	logger := slog.New(slog.Default().Handler())
	slog.SetDefault(logger)
	logger.Info("Starting custom agent service...")

	// open env file
	config.InitEnv(logger)

	r := gin.Default()

	redisClient := config.NewRedisClient()
	logger.Info("Connected to Redis")
	dynCfg := config.NewDynamicConfig(redisClient, logger, 1*time.Minute)

	queue := repository.NewRedisQueue(redisClient, "webhook_queue")
	webhookService := service.NewWebhookService(queue, logger)
	WebhookHandler := handler.NewWebhookHandler(webhookService)

	// start worker
	worker := worker.NewWebhookWorker(queue, logger, dynCfg)
	worker.Start(context.Background())

	r.POST("/webhook", WebhookHandler.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}
	r.Run(":" + port)
}
