package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/victor-divo/custom-agent-service/internal/config"
	"github.com/victor-divo/custom-agent-service/internal/handler"
)

func main() {
	config.InitEnv()

	r := gin.Default()

	r.POST("/webhook", handler.WebhookHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	r.Run(":" + port)
}
