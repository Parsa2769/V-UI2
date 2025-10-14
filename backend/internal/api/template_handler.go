package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/models"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type TemplateHandler struct {
	templateService *service.TemplateService
	auditService    *service.AuditService
	log             *zap.Logger
}

func NewTemplateHandler(templateService *service.TemplateService, auditService *service.AuditService, log *zap.Logger) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		auditService:    auditService,
		log:             log,
	}
}

// CreateTemplate creates a new configuration template
// @Summary Create template
// @Description Create a new configuration template
// @Tags templates
// @Accept json
// @Produce json
// @Param template body models.ConfigTemplate true "Template configuration"
// @Success 201 {object} models.ConfigTemplate
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /templates [post]
func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var template models.ConfigTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.CreatedBy = user.ID

	if err := h.templateService.CreateTemplate(&template); err != nil {
		h.log.Error("Failed to create template", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "template.create", c.ClientIP(), gin.H{
		"name": template.Name,
		"type": template.Type,
	})

	c.JSON(http.StatusCreated, template)
}

// GetTemplate gets a template by ID
// @Summary Get template
// @Description Get a configuration template by ID
// @Tags templates
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} models.ConfigTemplate
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /templates/{id} [get]
func (h *TemplateHandler) GetTemplate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	template, err := h.templateService.GetTemplate(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// ListTemplates lists all templates
// @Summary List templates
// @Description List all configuration templates
// @Tags templates
// @Produce json
// @Param type query string false "Template type filter"
// @Success 200 {array} models.ConfigTemplate
// @Security BearerAuth
// @Router /templates [get]
func (h *TemplateHandler) ListTemplates(c *gin.Context) {
	templateType := c.Query("type")

	templates, err := h.templateService.ListTemplates(templateType)
	if err != nil {
		h.log.Error("Failed to list templates", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// UpdateTemplate updates a template
// @Summary Update template
// @Description Update a configuration template
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param template body models.ConfigTemplate true "Updated template"
// @Success 200 {object} models.ConfigTemplate
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /templates/{id} [put]
func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var updates models.ConfigTemplate
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.templateService.UpdateTemplate(id, &updates); err != nil {
		h.log.Error("Failed to update template", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "template.update", c.ClientIP(), gin.H{
		"id": id.String(),
	})

	template, _ := h.templateService.GetTemplate(id)
	c.JSON(http.StatusOK, template)
}

// DeleteTemplate deletes a template
// @Summary Delete template
// @Description Delete a configuration template
// @Tags templates
// @Param id path string true "Template ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /templates/{id} [delete]
func (h *TemplateHandler) DeleteTemplate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.templateService.DeleteTemplate(id); err != nil {
		h.log.Error("Failed to delete template", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "template.delete", c.ClientIP(), gin.H{
		"id": id.String(),
	})

	c.Status(http.StatusNoContent)
}

// GetDefaultTemplates gets built-in default templates
// @Summary Get default templates
// @Description Get built-in default configuration templates
// @Tags templates
// @Produce json
// @Success 200 {array} models.ConfigTemplate
// @Security BearerAuth
// @Router /templates/defaults [get]
func (h *TemplateHandler) GetDefaultTemplates(c *gin.Context) {
	templates := h.templateService.GetDefaultTemplates()
	c.JSON(http.StatusOK, templates)
}

// ApplyTemplate applies a template to create an inbound
// @Summary Apply template
// @Description Apply a template to create a new inbound
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param overrides body map[string]interface{} false "Configuration overrides"
// @Success 200 {object} models.Inbound
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /templates/{id}/apply [post]
func (h *TemplateHandler) ApplyTemplate(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var overrides map[string]interface{}
	if err := c.ShouldBindJSON(&overrides); err != nil {
		overrides = make(map[string]interface{})
	}

	inbound, err := h.templateService.ApplyTemplate(id, overrides)
	if err != nil {
		h.log.Error("Failed to apply template", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "template.apply", c.ClientIP(), gin.H{
		"template_id": id.String(),
	})

	c.JSON(http.StatusOK, inbound)
}
