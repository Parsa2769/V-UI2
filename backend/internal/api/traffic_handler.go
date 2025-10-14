package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type TrafficHandler struct {
	trafficService *service.TrafficService
	log            *zap.Logger
}

func NewTrafficHandler(trafficService *service.TrafficService, log *zap.Logger) *TrafficHandler {
	return &TrafficHandler{
		trafficService: trafficService,
		log:            log,
	}
}

func (h *TrafficHandler) GetClientTraffic(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	from, _ := time.Parse(time.RFC3339, c.DefaultQuery("from", time.Now().AddDate(0, 0, -7).Format(time.RFC3339)))
	to, _ := time.Parse(time.RFC3339, c.DefaultQuery("to", time.Now().Format(time.RFC3339)))

	logs, err := h.trafficService.GetClientTraffic(id, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

func (h *TrafficHandler) GetNodeTraffic(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node ID"})
		return
	}

	from, _ := time.Parse(time.RFC3339, c.DefaultQuery("from", time.Now().AddDate(0, 0, -7).Format(time.RFC3339)))
	to, _ := time.Parse(time.RFC3339, c.DefaultQuery("to", time.Now().Format(time.RFC3339)))

	logs, err := h.trafficService.GetNodeTraffic(id, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

