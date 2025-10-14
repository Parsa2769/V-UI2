package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"github.com/v-ui/backend/internal/xray/protocols"
	"go.uber.org/zap"
)

type InboundService struct {
	db     *db.Database
	log    *zap.Logger
	vmess  *protocols.VMess
	vless  *protocols.VLESS
	trojan *protocols.Trojan
	ss     *protocols.Shadowsocks
}

func NewInboundService(database *db.Database, log *zap.Logger) *InboundService {
	return &InboundService{
		db:     database,
		log:    log,
		vmess:  &protocols.VMess{},
		vless:  &protocols.VLESS{},
		trojan: &protocols.Trojan{},
		ss:     &protocols.Shadowsocks{},
	}
}

// CreateInbound creates a new inbound
func (s *InboundService) CreateInbound(inbound *models.Inbound) error {
	// Validate settings based on protocol
	var err error
	switch inbound.Protocol {
	case "vmess":
		err = s.vmess.ValidateSettings(&inbound.Settings)
	case "vless":
		err = s.vless.ValidateSettings(&inbound.Settings)
	case "trojan":
		err = s.trojan.ValidateSettings(&inbound.Settings)
	case "shadowsocks":
		err = s.ss.ValidateSettings(&inbound.Settings)
	default:
		return fmt.Errorf("unsupported protocol: %s", inbound.Protocol)
	}

	if err != nil {
		return fmt.Errorf("invalid settings: %w", err)
	}

	// Check if tag already exists
	var existing models.Inbound
	if err := s.db.Where("tag = ?", inbound.Tag).First(&existing).Error; err == nil {
		return fmt.Errorf("inbound with tag '%s' already exists", inbound.Tag)
	}

	// Check if port is already in use
	if err := s.db.Where("port = ? AND enable = ?", inbound.Port, true).First(&existing).Error; err == nil {
		return fmt.Errorf("port %d is already in use", inbound.Port)
	}

	if err := s.db.Create(inbound).Error; err != nil {
		s.log.Error("Failed to create inbound", zap.Error(err))
		return err
	}

	s.log.Info("Inbound created", zap.String("tag", inbound.Tag), zap.String("protocol", inbound.Protocol))
	return nil
}

// GetInbound gets an inbound by ID
func (s *InboundService) GetInbound(id uuid.UUID) (*models.Inbound, error) {
	var inbound models.Inbound
	if err := s.db.First(&inbound, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inbound, nil
}

// GetInboundByTag gets an inbound by tag
func (s *InboundService) GetInboundByTag(tag string) (*models.Inbound, error) {
	var inbound models.Inbound
	if err := s.db.First(&inbound, "tag = ?", tag).Error; err != nil {
		return nil, err
	}
	return &inbound, nil
}

// ListInbounds lists all inbounds with pagination
func (s *InboundService) ListInbounds(page, pageSize int, userID *uuid.UUID) ([]models.Inbound, int64, error) {
	var inbounds []models.Inbound
	var total int64

	query := s.db.Model(&models.Inbound{})
	
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&inbounds).Error; err != nil {
		return nil, 0, err
	}

	return inbounds, total, nil
}

// UpdateInbound updates an inbound
func (s *InboundService) UpdateInbound(id uuid.UUID, updates *models.Inbound) error {
	var inbound models.Inbound
	if err := s.db.First(&inbound, "id = ?", id).Error; err != nil {
		return err
	}

	// Validate settings if protocol changed
	if updates.Protocol != "" && updates.Protocol != inbound.Protocol {
		var err error
		switch updates.Protocol {
		case "vmess":
			err = s.vmess.ValidateSettings(&updates.Settings)
		case "vless":
			err = s.vless.ValidateSettings(&updates.Settings)
		case "trojan":
			err = s.trojan.ValidateSettings(&updates.Settings)
		case "shadowsocks":
			err = s.ss.ValidateSettings(&updates.Settings)
		}
		if err != nil {
			return fmt.Errorf("invalid settings: %w", err)
		}
	}

	if err := s.db.Model(&inbound).Updates(updates).Error; err != nil {
		s.log.Error("Failed to update inbound", zap.Error(err))
		return err
	}

	s.log.Info("Inbound updated", zap.String("id", id.String()))
	return nil
}

// DeleteInbound deletes an inbound
func (s *InboundService) DeleteInbound(id uuid.UUID) error {
	if err := s.db.Delete(&models.Inbound{}, "id = ?", id).Error; err != nil {
		s.log.Error("Failed to delete inbound", zap.Error(err))
		return err
	}

	s.log.Info("Inbound deleted", zap.String("id", id.String()))
	return nil
}

// AddClient adds a client to an inbound
func (s *InboundService) AddClient(inboundID uuid.UUID, client models.InboundClient) error {
	var inbound models.Inbound
	if err := s.db.First(&inbound, "id = ?", inboundID).Error; err != nil {
		return err
	}

	// Generate client ID based on protocol
	switch inbound.Protocol {
	case "vmess":
		client = s.vmess.GenerateClient(client.Email, client.ID)
	case "vless":
		client = s.vless.GenerateClient(client.Email, client.ID, client.Flow)
	case "trojan":
		client = s.trojan.GenerateClient(client.Email, client.ID)
	case "shadowsocks":
		client = s.ss.GenerateClient(client.Email, client.ID, inbound.Settings.Method)
	}

	// Add to clients list
	inbound.Settings.Clients = append(inbound.Settings.Clients, client)

	if err := s.db.Model(&inbound).Update("settings", inbound.Settings).Error; err != nil {
		return err
	}

	s.log.Info("Client added to inbound", 
		zap.String("inbound", inbound.Tag),
		zap.String("email", client.Email))
	return nil
}

// RemoveClient removes a client from an inbound
func (s *InboundService) RemoveClient(inboundID uuid.UUID, clientEmail string) error {
	var inbound models.Inbound
	if err := s.db.First(&inbound, "id = ?", inboundID).Error; err != nil {
		return err
	}

	// Find and remove client
	newClients := make([]models.InboundClient, 0)
	found := false
	for _, c := range inbound.Settings.Clients {
		if c.Email != clientEmail {
			newClients = append(newClients, c)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("client not found: %s", clientEmail)
	}

	inbound.Settings.Clients = newClients

	if err := s.db.Model(&inbound).Update("settings", inbound.Settings).Error; err != nil {
		return err
	}

	s.log.Info("Client removed from inbound",
		zap.String("inbound", inbound.Tag),
		zap.String("email", clientEmail))
	return nil
}

// GenerateShareLink generates a share link for a client
func (s *InboundService) GenerateShareLink(inboundID uuid.UUID, clientEmail, address string) (string, error) {
	var inbound models.Inbound
	if err := s.db.First(&inbound, "id = ?", inboundID).Error; err != nil {
		return "", err
	}

	// Find client
	var client *models.InboundClient
	for i := range inbound.Settings.Clients {
		if inbound.Settings.Clients[i].Email == clientEmail {
			client = &inbound.Settings.Clients[i]
			break
		}
	}

	if client == nil {
		return "", fmt.Errorf("client not found: %s", clientEmail)
	}

	// Generate link based on protocol
	var link string
	var err error

	switch inbound.Protocol {
	case "vmess":
		link, err = s.vmess.GenerateShareLink(&inbound, client, address)
	case "vless":
		link, err = s.vless.GenerateShareLink(&inbound, client, address)
	case "trojan":
		link, err = s.trojan.GenerateShareLink(&inbound, client, address)
	case "shadowsocks":
		link, err = s.ss.GenerateShareLink(&inbound, client, address)
	default:
		return "", fmt.Errorf("unsupported protocol: %s", inbound.Protocol)
	}

	return link, err
}

// ToggleInbound enables or disables an inbound
func (s *InboundService) ToggleInbound(id uuid.UUID, enable bool) error {
	if err := s.db.Model(&models.Inbound{}).Where("id = ?", id).Update("enable", enable).Error; err != nil {
		return err
	}

	s.log.Info("Inbound toggled", zap.String("id", id.String()), zap.Bool("enable", enable))
	return nil
}

// GetInboundStats gets traffic statistics for an inbound
func (s *InboundService) GetInboundStats(id uuid.UUID) (map[string]interface{}, error) {
	var inbound models.Inbound
	if err := s.db.First(&inbound, "id = ?", id).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"id":            inbound.ID,
		"tag":           inbound.Tag,
		"protocol":      inbound.Protocol,
		"port":          inbound.Port,
		"clients_count": len(inbound.Settings.Clients),
		"enable":        inbound.Enable,
	}

	return stats, nil
}
