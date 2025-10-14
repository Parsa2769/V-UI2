package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/auth"
	"github.com/v-ui/backend/internal/models"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type ClientHandler struct {
	clientService *service.ClientService
	auditService  *service.AuditService
	log           *zap.Logger
}

func NewClientHandler(clientService *service.ClientService, auditService *service.AuditService, log *zap.Logger) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
		auditService:  auditService,
		log:           log,
	}
}

func (h *ClientHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	filters := make(map[string]interface{})
	if protocol := c.Query("protocol"); protocol != "" {
		filters["protocol = ?"] = protocol
	}
	if enable := c.Query("enable"); enable != "" {
		filters["enable = ?"] = enable == "true"
	}

	clients, total, err := h.clientService.List(offset, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  clients,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *ClientHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	client, err := h.clientService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) Create(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client.CreatedBy = claims.UserID
	if err := h.clientService.Create(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "client_create", "client:"+client.ID.String(), "Client created: "+client.Email, c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusCreated, client)
}

func (h *ClientHandler) Update(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.clientService.Update(id, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "client_update", "client:"+id.String(), "Client updated", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"message": "client updated successfully"})
}

func (h *ClientHandler) Delete(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	if err := h.clientService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "client_delete", "client:"+id.String(), "Client deleted", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"message": "client deleted successfully"})
}

func (h *ClientHandler) GetTraffic(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	upload, download, err := h.clientService.GetTrafficStats(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_bytes":   upload,
		"download_bytes": download,
		"total_bytes":    upload + download,
	})
}

