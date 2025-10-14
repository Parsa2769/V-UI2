package protocols

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/v-ui/backend/internal/models"
)

// Trojan protocol handler
type Trojan struct{}

// GenerateClient generates a Trojan client configuration
func (t *Trojan) GenerateClient(email, password string) models.InboundClient {
	return models.InboundClient{
		ID:     password, // Trojan uses password instead of UUID
		Email:  email,
		Enable: true,
	}
}

// GenerateShareLink generates a trojan:// share link
func (t *Trojan) GenerateShareLink(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	// trojan://password@address:port?parameters#remark
	
	params := url.Values{}
	params.Set("type", inbound.StreamSettings.Network)
	
	if inbound.StreamSettings.Security != "" {
		params.Set("security", inbound.StreamSettings.Security)
	} else {
		params.Set("security", "tls")
	}

	// Network-specific parameters
	switch inbound.StreamSettings.Network {
	case "ws":
		if inbound.StreamSettings.WSSettings != nil {
			params.Set("path", inbound.StreamSettings.WSSettings.Path)
			if len(inbound.StreamSettings.WSSettings.Headers.Request.Headers) > 0 {
				if host, ok := inbound.StreamSettings.WSSettings.Headers.Request.Headers["Host"]; ok && len(host) > 0 {
					params.Set("host", host[0])
				}
			}
		}
	case "grpc":
		if inbound.StreamSettings.GRPCSettings != nil {
			params.Set("serviceName", inbound.StreamSettings.GRPCSettings.ServiceName)
			params.Set("mode", "gun")
		}
	}

	// TLS settings
	if inbound.StreamSettings.TLSSettings != nil {
		if inbound.StreamSettings.TLSSettings.ServerName != "" {
			params.Set("sni", inbound.StreamSettings.TLSSettings.ServerName)
		}
		if len(inbound.StreamSettings.TLSSettings.ALPN) > 0 {
			params.Set("alpn", strings.Join(inbound.StreamSettings.TLSSettings.ALPN, ","))
		}
	}

	remark := url.QueryEscape(client.Email)
	link := fmt.Sprintf("trojan://%s@%s:%d?%s#%s",
		client.ID,
		address,
		inbound.Port,
		params.Encode(),
		remark,
	)

	return link, nil
}

// GenerateQRData generates QR code data
func (t *Trojan) GenerateQRData(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	return t.GenerateShareLink(inbound, client, address)
}

// ValidateSettings validates Trojan settings
func (t *Trojan) ValidateSettings(settings *models.InboundSettings) error {
	if len(settings.Clients) == 0 {
		return fmt.Errorf("at least one client is required")
	}

	for _, client := range settings.Clients {
		if client.ID == "" {
			return fmt.Errorf("client password is required")
		}
		if client.Email == "" {
			return fmt.Errorf("client email is required")
		}
	}

	return nil
}

// GetDefaultSettings returns default Trojan settings
func (t *Trojan) GetDefaultSettings() models.InboundSettings {
	return models.InboundSettings{
		Clients: []models.InboundClient{},
	}
}
