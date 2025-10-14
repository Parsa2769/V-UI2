package protocols

import (
	"fmt"

	"github.com/v-ui/backend/internal/models"
)

// Dokodemo (Tunnel/Door) protocol handler
type Dokodemo struct{}

// GenerateClient generates a Dokodemo-door configuration
// Note: Dokodemo-door is typically used for port forwarding/tunneling
func (d *Dokodemo) GenerateClient(email string) models.InboundClient {
	return models.InboundClient{
		Email:  email,
		Enable: true,
	}
}

// ValidateSettings validates Dokodemo-door settings
func (d *Dokodemo) ValidateSettings(settings *models.InboundSettings) error {
	// Dokodemo-door requires target address and port
	if settings.Network == "" {
		return fmt.Errorf("network type is required (tcp, udp, tcp,udp)")
	}

	return nil
}

// GetDefaultSettings returns default Dokodemo-door settings
func (d *Dokodemo) GetDefaultSettings(targetAddress string, targetPort int) models.InboundSettings {
	return models.InboundSettings{
		Network: "tcp,udp",
		// Target address and port would be stored in a custom field
		// or in the inbound configuration
	}
}

// GetTunnelInfo returns tunnel configuration info
func (d *Dokodemo) GetTunnelInfo(inbound *models.Inbound) map[string]interface{} {
	return map[string]interface{}{
		"listen_port": inbound.Port,
		"protocol":    "dokodemo-door",
		"network":     inbound.Settings.Network,
		"description": "Port forwarding / Transparent proxy",
	}
}

// GenerateConfig generates Xray config for Dokodemo-door
func (d *Dokodemo) GenerateConfig(listenPort int, targetAddress string, targetPort int, network string) map[string]interface{} {
	return map[string]interface{}{
		"port":     listenPort,
		"protocol": "dokodemo-door",
		"settings": map[string]interface{}{
			"address": targetAddress,
			"port":    targetPort,
			"network": network,
		},
		"tag": fmt.Sprintf("tunnel-%d", listenPort),
	}
}

// Tunnel is an alias for Dokodemo
type Tunnel struct {
	*Dokodemo
}

// NewTunnel creates a new Tunnel handler
func NewTunnel() *Tunnel {
	return &Tunnel{
		Dokodemo: &Dokodemo{},
	}
}
