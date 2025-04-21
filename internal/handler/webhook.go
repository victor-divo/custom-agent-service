package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victor-divo/custom-agent-service/internal/model"
)

func WebhookHandler(c *gin.Context) {
	var payload model.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	// Process the payload as needed

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully"})
}
