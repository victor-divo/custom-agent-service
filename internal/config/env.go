package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func InitEnv(logger *slog.Logger) {
	err := godotenv.Load()
	if err != nil {
		logger.Info("Error loading .env file")
	}
}

func GetAppId() string {
	return os.Getenv("QISCUS_APP_ID")
}

func GetSecretKey() string {
	return os.Getenv("QISCUS_SECRET_KEY")
}

func GetBaseURL() string {
	return os.Getenv("QISCUS_BASE_URL")
}
