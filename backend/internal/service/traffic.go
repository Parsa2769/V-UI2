package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type TrafficService struct {
	db  *db.Database
	log *zap.Logger
}

func NewTrafficService(database *db.Database, log *zap.Logger) *TrafficService {
	return &TrafficService{
		db:  database,
		log: log,
	}
}

func (s *TrafficService) RecordTraffic(clientID, nodeID uuid.UUID, uploadBytes, downloadBytes int64) error {
	log := &models.TrafficLog{
		ClientID:      clientID,
		NodeID:        &nodeID,
		UploadBytes:   uploadBytes,
		DownloadBytes: downloadBytes,
		Timestamp:     time.Now(),
	}
	return s.db.Create(log).Error
}

func (s *TrafficService) GetClientTraffic(clientID uuid.UUID, from, to time.Time) ([]models.TrafficLog, error) {
	var logs []models.TrafficLog
	err := s.db.Where("client_id = ? AND timestamp BETWEEN ? AND ?", clientID, from, to).
		Order("timestamp DESC").
		Find(&logs).Error
	return logs, err
}

func (s *TrafficService) GetNodeTraffic(nodeID uuid.UUID, from, to time.Time) ([]models.TrafficLog, error) {
	var logs []models.TrafficLog
	err := s.db.Where("node_id = ? AND timestamp BETWEEN ? AND ?", nodeID, from, to).
		Order("timestamp DESC").
		Find(&logs).Error
	return logs, err
}

