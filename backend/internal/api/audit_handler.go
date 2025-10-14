package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type AuditHandler struct {
	auditService *service.AuditService
	log          *zap.Logger
}

func NewAuditHandler(auditService *service.AuditService, log *zap.Logger) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
		log:          log,
	}
}

func (h *AuditHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	filters := make(map[string]interface{})
	if action := c.Query("action"); action != "" {
		filters["action = ?"] = action
	}

	logs, total, err := h.auditService.List(offset, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

