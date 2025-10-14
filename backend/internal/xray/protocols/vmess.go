package protocols

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/models"
)

// VMess protocol handler
type VMess struct{}

// GenerateClient generates a VMess client configuration
func (v *VMess) GenerateClient(email, id string) models.InboundClient {
	clientID := id
	if clientID == "" {
		clientID = uuid.New().String()
	}

	return models.InboundClient{
		ID:     clientID,
		Email:  email,
		Enable: true,
	}
}

// GenerateShareLink generates a vmess:// share link
func (v *VMess) GenerateShareLink(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	vmessConfig := map[string]interface{}{
		"v":    "2",
		"ps":   client.Email,
		"add":  address,
		"port": inbound.Port,
		"id":   client.ID,
		"aid":  "0",
		"net":  "tcp",
		"type": "none",
		"host": "",
		"path": "",
		"tls":  "",
	}

	// Apply stream settings
	if inbound.StreamSettings.Network != "" {
		vmessConfig["net"] = inbound.StreamSettings.Network
	}

	if inbound.StreamSettings.Security != "" {
		vmessConfig["tls"] = inbound.StreamSettings.Security
	}

	// Network-specific settings
	switch inbound.StreamSettings.Network {
	case "ws":
		if inbound.StreamSettings.WSSettings != nil {
			vmessConfig["path"] = inbound.StreamSettings.WSSettings.Path
			if len(inbound.StreamSettings.WSSettings.Headers.Request.Headers) > 0 {
				if host, ok := inbound.StreamSettings.WSSettings.Headers.Request.Headers["Host"]; ok && len(host) > 0 {
					vmessConfig["host"] = host[0]
				}
			}
		}
	case "http", "h2":
		if inbound.StreamSettings.HTTPSettings != nil {
			vmessConfig["path"] = inbound.StreamSettings.HTTPSettings.Path
			if len(inbound.StreamSettings.HTTPSettings.Host) > 0 {
				vmessConfig["host"] = inbound.StreamSettings.HTTPSettings.Host[0]
			}
		}
	case "grpc":
		if inbound.StreamSettings.GRPCSettings != nil {
			vmessConfig["path"] = inbound.StreamSettings.GRPCSettings.ServiceName
			vmessConfig["type"] = "gun"
		}
	case "quic":
		if inbound.StreamSettings.QUICSettings != nil {
			vmessConfig["type"] = inbound.StreamSettings.QUICSettings.Header.Type
			vmessConfig["host"] = inbound.StreamSettings.QUICSettings.Security
			vmessConfig["path"] = inbound.StreamSettings.QUICSettings.Key
		}
	case "kcp":
		if inbound.StreamSettings.KCPSettings != nil {
			vmessConfig["type"] = inbound.StreamSettings.KCPSettings.Header.Type
		}
	}

	// TLS settings
	if inbound.StreamSettings.TLSSettings != nil {
		vmessConfig["sni"] = inbound.StreamSettings.TLSSettings.ServerName
		if len(inbound.StreamSettings.TLSSettings.ALPN) > 0 {
			vmessConfig["alpn"] = inbound.StreamSettings.TLSSettings.ALPN[0]
		}
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(vmessConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vmess config: %w", err)
	}

	// Base64 encode
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return "vmess://" + encoded, nil
}

// GenerateQRData generates QR code data
func (v *VMess) GenerateQRData(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	return v.GenerateShareLink(inbound, client, address)
}

// ValidateSettings validates VMess settings
func (v *VMess) ValidateSettings(settings *models.InboundSettings) error {
	if len(settings.Clients) == 0 {
		return fmt.Errorf("at least one client is required")
	}

	for _, client := range settings.Clients {
		if client.ID == "" {
			return fmt.Errorf("client ID is required")
		}
		if client.Email == "" {
			return fmt.Errorf("client email is required")
		}
		// Validate UUID format
		if _, err := uuid.Parse(client.ID); err != nil {
			return fmt.Errorf("invalid client ID format: %w", err)
		}
	}

	return nil
}

// GetDefaultSettings returns default VMess settings
func (v *VMess) GetDefaultSettings() models.InboundSettings {
	return models.InboundSettings{
		Clients: []models.InboundClient{},
	}
}
