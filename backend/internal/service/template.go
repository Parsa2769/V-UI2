package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type TemplateService struct {
	db  *db.Database
	log *zap.Logger
}

func NewTemplateService(database *db.Database, log *zap.Logger) *TemplateService {
	return &TemplateService{
		db:  database,
		log: log,
	}
}

// CreateTemplate creates a new configuration template
func (s *TemplateService) CreateTemplate(template *models.ConfigTemplate) error {
	// Validate JSON
	var test interface{}
	if err := json.Unmarshal([]byte(template.Config), &test); err != nil {
		return fmt.Errorf("invalid JSON config: %w", err)
	}

	// Check if name already exists
	var existing models.ConfigTemplate
	if err := s.db.Where("name = ?", template.Name).First(&existing).Error; err == nil {
		return fmt.Errorf("template with name '%s' already exists", template.Name)
	}

	if err := s.db.Create(template).Error; err != nil {
		s.log.Error("Failed to create template", zap.Error(err))
		return err
	}

	s.log.Info("Template created", zap.String("name", template.Name))
	return nil
}

// GetTemplate gets a template by ID
func (s *TemplateService) GetTemplate(id uuid.UUID) (*models.ConfigTemplate, error) {
	var template models.ConfigTemplate
	if err := s.db.First(&template, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// GetTemplateByName gets a template by name
func (s *TemplateService) GetTemplateByName(name string) (*models.ConfigTemplate, error) {
	var template models.ConfigTemplate
	if err := s.db.First(&template, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// ListTemplates lists all templates
func (s *TemplateService) ListTemplates(templateType string) ([]models.ConfigTemplate, error) {
	var templates []models.ConfigTemplate
	
	query := s.db.Model(&models.ConfigTemplate{})
	
	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}

	if err := query.Order("name").Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

// UpdateTemplate updates a template
func (s *TemplateService) UpdateTemplate(id uuid.UUID, updates *models.ConfigTemplate) error {
	var template models.ConfigTemplate
	if err := s.db.First(&template, "id = ?", id).Error; err != nil {
		return err
	}

	// Validate JSON if config is being updated
	if updates.Config != "" {
		var test interface{}
		if err := json.Unmarshal([]byte(updates.Config), &test); err != nil {
			return fmt.Errorf("invalid JSON config: %w", err)
		}
	}

	if err := s.db.Model(&template).Updates(updates).Error; err != nil {
		s.log.Error("Failed to update template", zap.Error(err))
		return err
	}

	s.log.Info("Template updated", zap.String("id", id.String()))
	return nil
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(id uuid.UUID) error {
	if err := s.db.Delete(&models.ConfigTemplate{}, "id = ?", id).Error; err != nil {
		s.log.Error("Failed to delete template", zap.Error(err))
		return err
	}

	s.log.Info("Template deleted", zap.String("id", id.String()))
	return nil
}

// ApplyTemplate applies a template to create an inbound
func (s *TemplateService) ApplyTemplate(templateID uuid.UUID, overrides map[string]interface{}) (*models.Inbound, error) {
	// Get template
	template, err := s.GetTemplate(templateID)
	if err != nil {
		return nil, err
	}

	// Parse template config
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(template.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	// Apply overrides
	for key, value := range overrides {
		config[key] = value
	}

	// Create inbound from template
	inbound := &models.Inbound{
		Protocol: template.Type,
		Remark:   fmt.Sprintf("From template: %s", template.Name),
	}

	// This would need more sophisticated mapping
	// For now, return the structure
	s.log.Info("Template applied", zap.String("template", template.Name))
	return inbound, nil
}

// GetDefaultTemplates returns built-in default templates
func (s *TemplateService) GetDefaultTemplates() []models.ConfigTemplate {
	return []models.ConfigTemplate{
		{
			Name:        "VLESS-Reality-Vision",
			Type:        "vless",
			Description: "VLESS with Reality and Vision flow for maximum security",
			Config: `{
				"protocol": "vless",
				"settings": {
					"clients": [],
					"decryption": "none"
				},
				"streamSettings": {
					"network": "tcp",
					"security": "reality",
					"realitySettings": {
						"show": false,
						"dest": "www.google.com:443",
						"serverNames": ["www.google.com"],
						"privateKey": "GENERATE_NEW",
						"shortIds": [""]
					}
				}
			}`,
		},
		{
			Name:        "VMess-WebSocket-TLS",
			Type:        "vmess",
			Description: "VMess over WebSocket with TLS - Compatible with CDN",
			Config: `{
				"protocol": "vmess",
				"settings": {
					"clients": []
				},
				"streamSettings": {
					"network": "ws",
					"security": "tls",
					"wsSettings": {
						"path": "/api",
						"headers": {}
					},
					"tlsSettings": {
						"serverName": "example.com",
						"certificates": []
					}
				}
			}`,
		},
		{
			Name:        "Trojan-gRPC",
			Type:        "trojan",
			Description: "Trojan over gRPC for better performance",
			Config: `{
				"protocol": "trojan",
				"settings": {
					"clients": []
				},
				"streamSettings": {
					"network": "grpc",
					"security": "tls",
					"grpcSettings": {
						"serviceName": "TrojanService",
						"multiMode": false
					}
				}
			}`,
		},
		{
			Name:        "Shadowsocks-2022",
			Type:        "shadowsocks",
			Description: "Shadowsocks with 2022 encryption - Fastest performance",
			Config: `{
				"protocol": "shadowsocks",
				"settings": {
					"method": "2022-blake3-aes-256-gcm",
					"password": "GENERATE_NEW",
					"network": "tcp,udp"
				}
			}`,
		},
	}
}

// CreateDefaultTemplates creates the default templates in database
func (s *TemplateService) CreateDefaultTemplates() error {
	defaults := s.GetDefaultTemplates()
	
	for _, template := range defaults {
		// Check if already exists
		var existing models.ConfigTemplate
		if err := s.db.Where("name = ?", template.Name).First(&existing).Error; err == nil {
			continue // Already exists
		}

		if err := s.db.Create(&template).Error; err != nil {
			s.log.Warn("Failed to create default template", 
				zap.String("name", template.Name),
				zap.Error(err))
			continue
		}

		s.log.Info("Default template created", zap.String("name", template.Name))
	}

	return nil
}
