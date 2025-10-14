package service

import (
	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type ClientService struct {
	db  *db.Database
	log *zap.Logger
}

func NewClientService(database *db.Database, log *zap.Logger) *ClientService {
	return &ClientService{
		db:  database,
		log: log,
	}
}

func (s *ClientService) Create(client *models.Client) error {
	return s.db.Create(client).Error
}

func (s *ClientService) GetByID(id uuid.UUID) (*models.Client, error) {
	var client models.Client
	if err := s.db.Preload("Node").First(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (s *ClientService) List(offset, limit int, filters map[string]interface{}) ([]models.Client, int64, error) {
	var clients []models.Client
	var total int64

	query := s.db.Model(&models.Client{})
	for k, v := range filters {
		query = query.Where(k, v)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Node").Offset(offset).Limit(limit).Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

func (s *ClientService) Update(id uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&models.Client{}).Where("id = ?", id).Updates(updates).Error
}

func (s *ClientService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Client{}, id).Error
}

func (s *ClientService) GetTrafficStats(clientID uuid.UUID) (int64, int64, error) {
	var client models.Client
	if err := s.db.First(&client, clientID).Error; err != nil {
		return 0, 0, err
	}
	return client.UploadBytes, client.DownloadBytes, nil
}

