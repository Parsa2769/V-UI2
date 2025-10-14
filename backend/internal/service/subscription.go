package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type SubscriptionService struct {
	db             *db.Database
	log            *zap.Logger
	inboundService *InboundService
}

func NewSubscriptionService(database *db.Database, log *zap.Logger, inboundSvc *InboundService) *SubscriptionService {
	return &SubscriptionService{
		db:             database,
		log:            log,
		inboundService: inboundSvc,
	}
}

// GenerateSubscriptionToken generates a unique subscription token for a user
func (s *SubscriptionService) GenerateSubscriptionToken(userID uuid.UUID) (string, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return "", err
	}

	// Generate token (use user ID + random string)
	token := base64.URLEncoding.EncodeToString([]byte(userID.String()))
	return token, nil
}

// GetSubscriptionContent generates subscription content for a user
func (s *SubscriptionService) GetSubscriptionContent(token, address string) (string, error) {
	// Decode token to get user ID
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(string(decoded))
	if err != nil {
		return "", fmt.Errorf("invalid token format")
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Get all inbounds for this user
	inbounds, _, err := s.inboundService.ListInbounds(1, 1000, &userID)
	if err != nil {
		return "", err
	}

	// Generate links for all clients
	var links []string
	for _, inbound := range inbounds {
		if !inbound.Enable {
			continue
		}

		for _, client := range inbound.Settings.Clients {
			if !client.Enable {
				continue
			}

			link, err := s.inboundService.GenerateShareLink(inbound.ID, client.Email, address)
			if err != nil {
				s.log.Warn("Failed to generate link", 
					zap.String("inbound", inbound.Tag),
					zap.String("client", client.Email),
					zap.Error(err))
				continue
			}

			links = append(links, link)
		}
	}

	// Join all links
	content := strings.Join(links, "\n")

	// Base64 encode
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	return encoded, nil
}

// GetUserClients gets all clients for a user across all inbounds
func (s *SubscriptionService) GetUserClients(userID uuid.UUID) ([]map[string]interface{}, error) {
	inbounds, _, err := s.inboundService.ListInbounds(1, 1000, &userID)
	if err != nil {
		return nil, err
	}

	var clients []map[string]interface{}
	for _, inbound := range inbounds {
		for _, client := range inbound.Settings.Clients {
			clientInfo := map[string]interface{}{
				"inbound_id":   inbound.ID,
				"inbound_tag":  inbound.Tag,
				"protocol":     inbound.Protocol,
				"port":         inbound.Port,
				"email":        client.Email,
				"id":           client.ID,
				"enable":       client.Enable,
				"limit_ip":     client.LimitIP,
				"total_gb":     client.TotalGB,
				"expiry_time":  client.ExpiryTime,
			}
			clients = append(clients, clientInfo)
		}
	}

	return clients, nil
}
