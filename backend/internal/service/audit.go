package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type AuditService struct {
	db  *db.Database
	log *zap.Logger
}

func NewAuditService(database *db.Database, log *zap.Logger) *AuditService {
	return &AuditService{
		db:  database,
		log: log,
	}
}

func (s *AuditService) Log(actorID uuid.UUID, action, target, detail, ipAddress, userAgent string) error {
	auditLog := &models.AuditLog{
		ActorID:   actorID,
		Action:    action,
		Target:    target,
		Detail:    detail,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}
	return s.db.Create(auditLog).Error
}

func (s *AuditService) List(offset, limit int, filters map[string]interface{}) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := s.db.Model(&models.AuditLog{})
	for k, v := range filters {
		query = query.Where(k, v)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Actor").Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

