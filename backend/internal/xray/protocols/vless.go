package protocols

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/models"
)

// VLESS protocol handler
type VLESS struct{}

// GenerateClient generates a VLESS client configuration
func (v *VLESS) GenerateClient(email, id, flow string) models.InboundClient {
	clientID := id
	if clientID == "" {
		clientID = uuid.New().String()
	}

	return models.InboundClient{
		ID:     clientID,
		Email:  email,
		Flow:   flow, // "", "xtls-rprx-vision", "xtls-rprx-direct"
		Enable: true,
	}
}

// GenerateShareLink generates a vless:// share link
func (v *VLESS) GenerateShareLink(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	// vless://UUID@address:port?parameters#remark
	
	params := url.Values{}
	params.Set("type", inbound.StreamSettings.Network)
	
	// Security
	if inbound.StreamSettings.Security != "" {
		params.Set("security", inbound.StreamSettings.Security)
	} else {
		params.Set("security", "none")
	}

	// Flow for XTLS
	if client.Flow != "" {
		params.Set("flow", client.Flow)
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
	case "http", "h2":
		if inbound.StreamSettings.HTTPSettings != nil {
			params.Set("path", inbound.StreamSettings.HTTPSettings.Path)
			if len(inbound.StreamSettings.HTTPSettings.Host) > 0 {
				params.Set("host", strings.Join(inbound.StreamSettings.HTTPSettings.Host, ","))
			}
		}
	case "grpc":
		if inbound.StreamSettings.GRPCSettings != nil {
			params.Set("serviceName", inbound.StreamSettings.GRPCSettings.ServiceName)
			params.Set("mode", "gun")
		}
	case "quic":
		if inbound.StreamSettings.QUICSettings != nil {
			params.Set("quicSecurity", inbound.StreamSettings.QUICSettings.Security)
			params.Set("key", inbound.StreamSettings.QUICSettings.Key)
			params.Set("headerType", inbound.StreamSettings.QUICSettings.Header.Type)
		}
	case "kcp":
		if inbound.StreamSettings.KCPSettings != nil {
			params.Set("headerType", inbound.StreamSettings.KCPSettings.Header.Type)
			params.Set("seed", "")
		}
	case "tcp":
		if inbound.StreamSettings.TCPSettings != nil && inbound.StreamSettings.TCPSettings.Header.Type != "" {
			params.Set("headerType", inbound.StreamSettings.TCPSettings.Header.Type)
		}
	}

	// TLS/Reality settings
	if inbound.StreamSettings.Security == "tls" && inbound.StreamSettings.TLSSettings != nil {
		if inbound.StreamSettings.TLSSettings.ServerName != "" {
			params.Set("sni", inbound.StreamSettings.TLSSettings.ServerName)
		}
		if len(inbound.StreamSettings.TLSSettings.ALPN) > 0 {
			params.Set("alpn", strings.Join(inbound.StreamSettings.TLSSettings.ALPN, ","))
		}
	}

	if inbound.StreamSettings.Security == "reality" && inbound.StreamSettings.RealitySettings != nil {
		if len(inbound.StreamSettings.RealitySettings.ServerNames) > 0 {
			params.Set("sni", inbound.StreamSettings.RealitySettings.ServerNames[0])
		}
		params.Set("pbk", inbound.StreamSettings.RealitySettings.Settings.PublicKey)
		params.Set("fp", inbound.StreamSettings.RealitySettings.Settings.Fingerprint)
		if len(inbound.StreamSettings.RealitySettings.ShortIds) > 0 {
			params.Set("sid", inbound.StreamSettings.RealitySettings.ShortIds[0])
		}
		if inbound.StreamSettings.RealitySettings.Settings.SpiderX != "" {
			params.Set("spx", inbound.StreamSettings.RealitySettings.Settings.SpiderX)
		}
	}

	// Build URL
	remark := url.QueryEscape(client.Email)
	link := fmt.Sprintf("vless://%s@%s:%d?%s#%s",
		client.ID,
		address,
		inbound.Port,
		params.Encode(),
		remark,
	)

	return link, nil
}

// GenerateQRData generates QR code data
func (v *VLESS) GenerateQRData(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	return v.GenerateShareLink(inbound, client, address)
}

// ValidateSettings validates VLESS settings
func (v *VLESS) ValidateSettings(settings *models.InboundSettings) error {
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
		// Validate flow
		if client.Flow != "" {
			validFlows := []string{"xtls-rprx-vision", "xtls-rprx-direct"}
			valid := false
			for _, f := range validFlows {
				if client.Flow == f {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid flow: %s", client.Flow)
			}
		}
	}

	// Validate decryption
	if settings.Decryption != "" && settings.Decryption != "none" {
		return fmt.Errorf("invalid decryption: %s (must be 'none' or empty)", settings.Decryption)
	}

	return nil
}

// GetDefaultSettings returns default VLESS settings
func (v *VLESS) GetDefaultSettings() models.InboundSettings {
	return models.InboundSettings{
		Clients:    []models.InboundClient{},
		Decryption: "none",
	}
}
