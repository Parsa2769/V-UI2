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

type InboundHandler struct {
	inboundService *service.InboundService
	auditService   *service.AuditService
	log            *zap.Logger
}

func NewInboundHandler(inboundService *service.InboundService, auditService *service.AuditService, log *zap.Logger) *InboundHandler {
	return &InboundHandler{
		inboundService: inboundService,
		auditService:   auditService,
		log:            log,
	}
}

// CreateInbound creates a new inbound
// @Summary Create inbound
// @Description Create a new inbound configuration
// @Tags inbounds
// @Accept json
// @Produce json
// @Param inbound body models.Inbound true "Inbound configuration"
// @Success 201 {object} models.Inbound
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds [post]
func (h *InboundHandler) CreateInbound(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var inbound models.Inbound
	if err := c.ShouldBindJSON(&inbound); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inbound.UserID = user.ID

	if err := h.inboundService.CreateInbound(&inbound); err != nil {
		h.log.Error("Failed to create inbound", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "inbound.create", c.ClientIP(), gin.H{
		"tag":      inbound.Tag,
		"protocol": inbound.Protocol,
		"port":     inbound.Port,
	})

	c.JSON(http.StatusCreated, inbound)
}

// GetInbound gets an inbound by ID
// @Summary Get inbound
// @Description Get an inbound configuration by ID
// @Tags inbounds
// @Produce json
// @Param id path string true "Inbound ID"
// @Success 200 {object} models.Inbound
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id} [get]
func (h *InboundHandler) GetInbound(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	inbound, err := h.inboundService.GetInbound(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "inbound not found"})
		return
	}

	c.JSON(http.StatusOK, inbound)
}

// ListInbounds lists all inbounds
// @Summary List inbounds
// @Description List all inbound configurations with pagination
// @Tags inbounds
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds [get]
func (h *InboundHandler) ListInbounds(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var userID *uuid.UUID
	// Non-admin users can only see their own inbounds
	if user.Role != auth.RoleAdmin {
		userID = &user.ID
	}

	inbounds, total, err := h.inboundService.ListInbounds(page, pageSize, userID)
	if err != nil {
		h.log.Error("Failed to list inbounds", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch inbounds"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      inbounds,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// UpdateInbound updates an inbound
// @Summary Update inbound
// @Description Update an inbound configuration
// @Tags inbounds
// @Accept json
// @Produce json
// @Param id path string true "Inbound ID"
// @Param inbound body models.Inbound true "Updated inbound configuration"
// @Success 200 {object} models.Inbound
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id} [put]
func (h *InboundHandler) UpdateInbound(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var updates models.Inbound
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.inboundService.UpdateInbound(id, &updates); err != nil {
		h.log.Error("Failed to update inbound", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "inbound.update", c.ClientIP(), gin.H{
		"id": id.String(),
	})

	inbound, _ := h.inboundService.GetInbound(id)
	c.JSON(http.StatusOK, inbound)
}

// DeleteInbound deletes an inbound
// @Summary Delete inbound
// @Description Delete an inbound configuration
// @Tags inbounds
// @Param id path string true "Inbound ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id} [delete]
func (h *InboundHandler) DeleteInbound(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.inboundService.DeleteInbound(id); err != nil {
		h.log.Error("Failed to delete inbound", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "inbound.delete", c.ClientIP(), gin.H{
		"id": id.String(),
	})

	c.Status(http.StatusNoContent)
}

// AddClient adds a client to an inbound
// @Summary Add client
// @Description Add a client to an inbound
// @Tags inbounds
// @Accept json
// @Produce json
// @Param id path string true "Inbound ID"
// @Param client body models.InboundClient true "Client configuration"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id}/clients [post]
func (h *InboundHandler) AddClient(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var client models.InboundClient
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.inboundService.AddClient(id, client); err != nil {
		h.log.Error("Failed to add client", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "inbound.add_client", c.ClientIP(), gin.H{
		"inbound_id": id.String(),
		"email":      client.Email,
	})

	c.JSON(http.StatusOK, gin.H{"message": "client added successfully"})
}

// RemoveClient removes a client from an inbound
// @Summary Remove client
// @Description Remove a client from an inbound
// @Tags inbounds
// @Param id path string true "Inbound ID"
// @Param email path string true "Client email"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id}/clients/{email} [delete]
func (h *InboundHandler) RemoveClient(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	email := c.Param("email")

	if err := h.inboundService.RemoveClient(id, email); err != nil {
		h.log.Error("Failed to remove client", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "inbound.remove_client", c.ClientIP(), gin.H{
		"inbound_id": id.String(),
		"email":      email,
	})

	c.Status(http.StatusNoContent)
}

// GetShareLink gets a share link for a client
// @Summary Get share link
// @Description Get a share link for a specific client
// @Tags inbounds
// @Produce json
// @Param id path string true "Inbound ID"
// @Param email path string true "Client email"
// @Param address query string false "Server address"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id}/clients/{email}/link [get]
func (h *InboundHandler) GetShareLink(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	email := c.Param("email")
	address := c.DefaultQuery("address", c.Request.Host)

	link, err := h.inboundService.GenerateShareLink(id, email, address)
	if err != nil {
		h.log.Error("Failed to generate share link", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link": link,
	})
}

// ToggleInbound enables or disables an inbound
// @Summary Toggle inbound
// @Description Enable or disable an inbound
// @Tags inbounds
// @Accept json
// @Produce json
// @Param id path string true "Inbound ID"
// @Param body body map[string]bool true "Enable status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id}/toggle [post]
func (h *InboundHandler) ToggleInbound(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var body struct {
		Enable bool `json:"enable"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.inboundService.ToggleInbound(id, body.Enable); err != nil {
		h.log.Error("Failed to toggle inbound", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "inbound.toggle", c.ClientIP(), gin.H{
		"id":     id.String(),
		"enable": body.Enable,
	})

	c.JSON(http.StatusOK, gin.H{"message": "inbound toggled successfully"})
}

// GetInboundStats gets statistics for an inbound
// @Summary Get inbound stats
// @Description Get traffic statistics for an inbound
// @Tags inbounds
// @Produce json
// @Param id path string true "Inbound ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /inbounds/{id}/stats [get]
func (h *InboundHandler) GetInboundStats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	stats, err := h.inboundService.GetInboundStats(id)
	if err != nil {
		h.log.Error("Failed to get inbound stats", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
