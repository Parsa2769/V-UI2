package service

import (
	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type NodeService struct {
	db  *db.Database
	log *zap.Logger
}

func NewNodeService(database *db.Database, log *zap.Logger) *NodeService {
	return &NodeService{
		db:  database,
		log: log,
	}
}

func (s *NodeService) Create(node *models.Node) error {
	return s.db.Create(node).Error
}

func (s *NodeService) GetByID(id uuid.UUID) (*models.Node, error) {
	var node models.Node
	if err := s.db.First(&node, id).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *NodeService) List(offset, limit int) ([]models.Node, int64, error) {
	var nodes []models.Node
	var total int64

	if err := s.db.Model(&models.Node{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Offset(offset).Limit(limit).Find(&nodes).Error; err != nil {
		return nil, 0, err
	}

	return nodes, total, nil
}

func (s *NodeService) Update(id uuid.UUID, updates map[string]interface{}) error {
	return s.db.Model(&models.Node{}).Where("id = ?", id).Updates(updates).Error
}

func (s *NodeService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Node{}, id).Error
}

