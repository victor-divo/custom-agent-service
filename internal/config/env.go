package config

import (
	"log/slog"
	"os"
	"strconv"

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

func GetDefaultMaxAgentChats() int {
	defaultMaxChats, err := strconv.Atoi(os.Getenv("DEFAULT_MAX_AGENT_CHATS"))
	if err != nil {
		return 2
	}
	return defaultMaxChats
}
