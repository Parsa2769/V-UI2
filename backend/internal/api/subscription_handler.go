package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
	log                 *zap.Logger
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService, log *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		log:                 log,
	}
}

// GetSubscription handles subscription requests
// @Summary Get subscription
// @Description Get subscription content for a user
// @Tags subscription
// @Produce plain
// @Param token path string true "Subscription token"
// @Success 200 {string} string "Base64 encoded subscription content"
// @Failure 400 {object} map[string]interface{}
// @Router /sub/{token} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	// Get server address from request or config
	address := c.Request.Host

	content, err := h.subscriptionService.GetSubscriptionContent(token, address)
	if err != nil {
		h.log.Error("Failed to get subscription", 
			zap.String("token", token),
			zap.Error(err))
		c.String(http.StatusBadRequest, "Invalid subscription token")
		return
	}

	// Set headers for subscription
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=subscription.txt")
	c.Header("Profile-Update-Interval", "24")
	c.Header("Subscription-Userinfo", "upload=0; download=0; total=0; expire=0")

	c.String(http.StatusOK, content)
}

// GenerateToken generates a subscription token for a user
// @Summary Generate subscription token
// @Description Generate a subscription token for the current user
// @Tags subscription
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /subscription/token [get]
func (h *SubscriptionHandler) GenerateToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid := userID.(uuid.UUID)

	token, err := h.subscriptionService.GenerateSubscriptionToken(uid)
	if err != nil {
		h.log.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Build subscription URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	subscriptionURL := scheme + "://" + c.Request.Host + "/sub/" + token

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"url":   subscriptionURL,
	})
}

// GetUserClients gets all clients for the current user
// @Summary Get user clients
// @Description Get all clients for the current user
// @Tags subscription
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Security BearerAuth
// @Router /subscription/clients [get]
func (h *SubscriptionHandler) GetUserClients(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid := userID.(uuid.UUID)

	clients, err := h.subscriptionService.GetUserClients(uid)
	if err != nil {
		h.log.Error("Failed to get clients", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch clients"})
		return
	}

	c.JSON(http.StatusOK, clients)
}
