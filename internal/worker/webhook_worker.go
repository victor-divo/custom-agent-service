package worker

import (
	"context"
	"log/slog"
	"os"

	"github.com/victor-divo/custom-agent-service/internal/model"
	"github.com/victor-divo/custom-agent-service/internal/repository"
)

type WebhookWorker struct {
	Queue  repository.Queue
	logger *slog.Logger
}

func NewWebhookWorker(queue repository.Queue, logger *slog.Logger) *WebhookWorker {
	return &WebhookWorker{
		Queue:  queue,
		logger: logger,
	}
}

func (w *WebhookWorker) Start(ctx context.Context) {

	go func() {
		for {
			payload, err := w.Queue.Dequeue(ctx)
			if err != nil {
				w.logger.Error("Error dequeuing payload:", "error", err)
				continue
			}

			if payload == nil {
				continue
			}

			w.processPayload(payload)
		}
	}()
}

func (w *WebhookWorker) processPayload(payload *model.WebhookPayload) {
	// Contoh log, ganti sesuai logika kamu
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("processing data", "payload", payload)
}
