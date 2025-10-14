package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/v-ui/backend/internal/models"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type BackupHandler struct {
	backupService *service.BackupService
	auditService  *service.AuditService
	log           *zap.Logger
}

func NewBackupHandler(backupService *service.BackupService, auditService *service.AuditService, log *zap.Logger) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
		auditService:  auditService,
		log:           log,
	}
}

// CreateBackup creates a full database backup
// @Summary Create backup
// @Description Create a full database backup
// @Tags backup
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /backup/create [post]
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	filename, err := h.backupService.CreateBackup()
	if err != nil {
		h.log.Error("Failed to create backup", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create backup"})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "backup.create", c.ClientIP(), gin.H{
		"filename": filename,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":  "backup created successfully",
		"filename": filename,
	})
}

// ListBackups lists all available backups
// @Summary List backups
// @Description List all available backup files
// @Tags backup
// @Produce json
// @Success 200 {array} service.BackupInfo
// @Security BearerAuth
// @Router /backup/list [get]
func (h *BackupHandler) ListBackups(c *gin.Context) {
	backups, err := h.backupService.ListBackups()
	if err != nil {
		h.log.Error("Failed to list backups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list backups"})
		return
	}

	c.JSON(http.StatusOK, backups)
}

// RestoreBackup restores database from a backup file
// @Summary Restore backup
// @Description Restore database from a backup file
// @Tags backup
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Restore options"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /backup/restore [post]
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var body struct {
		Filename string                    `json:"filename" binding:"required"`
		Options  service.ImportOptions     `json:"options"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.backupService.RestoreFromFile(body.Filename, body.Options); err != nil {
		h.log.Error("Failed to restore backup", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "backup.restore", c.ClientIP(), gin.H{
		"filename": body.Filename,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "backup restored successfully",
	})
}

// DeleteBackup deletes a backup file
// @Summary Delete backup
// @Description Delete a backup file
// @Tags backup
// @Param filename path string true "Backup filename"
// @Success 204
// @Security BearerAuth
// @Router /backup/{filename} [delete]
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	filename := c.Param("filename")

	if err := h.backupService.DeleteBackup(filename); err != nil {
		h.log.Error("Failed to delete backup", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Audit log
	h.auditService.Log(user.ID, "backup.delete", c.ClientIP(), gin.H{
		"filename": filename,
	})

	c.Status(http.StatusNoContent)
}
