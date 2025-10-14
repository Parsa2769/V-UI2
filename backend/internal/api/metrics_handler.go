package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v-ui/backend/internal/models"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type MetricsHandler struct {
	services *service.Services
	log      *zap.Logger
}

func NewMetricsHandler(services *service.Services, log *zap.Logger) *MetricsHandler {
	return &MetricsHandler{
		services: services,
		log:      log,
	}
}

func (h *MetricsHandler) GetSummary(c *gin.Context) {
	// Get total users
	users, totalUsers, _ := h.services.User.List(0, 1)
	_ = users

	// Get total clients
	clients, totalClients, _ := h.services.Client.List(0, 1, nil)
	_ = clients

	// Get active clients
	activeClients, _, _ := h.services.Client.List(0, 1, map[string]interface{}{"enable = ?": true})

	// Get total nodes
	nodes, totalNodes, _ := h.services.Node.List(0, 100)

	// Count online nodes
	onlineNodes := 0
	for _, node := range nodes {
		if node.Status == models.NodeStatusOnline {
			onlineNodes++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users": gin.H{
			"total": totalUsers,
		},
		"clients": gin.H{
			"total":  totalClients,
			"active": len(activeClients),
		},
		"nodes": gin.H{
			"total":  totalNodes,
			"online": onlineNodes,
		},
	})
}

