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

type NodeHandler struct {
	nodeService  *service.NodeService
	auditService *service.AuditService
	log          *zap.Logger
}

func NewNodeHandler(nodeService *service.NodeService, auditService *service.AuditService, log *zap.Logger) *NodeHandler {
	return &NodeHandler{
		nodeService:  nodeService,
		auditService: auditService,
		log:          log,
	}
}

func (h *NodeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	nodes, total, err := h.nodeService.List(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  nodes,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *NodeHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node ID"})
		return
	}

	node, err := h.nodeService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	c.JSON(http.StatusOK, node)
}

func (h *NodeHandler) Create(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	var node models.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nodeService.Create(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "node_create", "node:"+node.ID.String(), "Node created: "+node.Name, c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusCreated, node)
}

func (h *NodeHandler) Update(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.nodeService.Update(id, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "node_update", "node:"+id.String(), "Node updated", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"message": "node updated successfully"})
}

func (h *NodeHandler) Delete(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node ID"})
		return
	}

	if err := h.nodeService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "node_delete", "node:"+id.String(), "Node deleted", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"message": "node deleted successfully"})
}

