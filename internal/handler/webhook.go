package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victor-divo/custom-agent-service/internal/model"
	"github.com/victor-divo/custom-agent-service/internal/service"
)

type WebhookHandler struct {
	Service *service.WebhookService
}

func NewWebhookHandler(svc *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		Service: svc,
	}
}

func (h *WebhookHandler) Handle(c *gin.Context) {
	var payload model.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := h.Service.HandleWebhook(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully"})
}

func (h *WebhookHandler) Resolve(c *gin.Context) {
	var payload model.ResolvePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}
	if err := h.Service.HandleResolvedChat(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received successfully"})
}
