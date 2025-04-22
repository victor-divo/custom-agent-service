package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func InitEnv(logger *slog.Logger) {
	err := godotenv.Load()
	if err != nil {
		logger.Info("Error loading .env file")
	}
}
