package protocols

import (
	"fmt"

	"github.com/v-ui/backend/internal/models"
)

// HTTP protocol handler (HTTP/SOCKS proxy)
type HTTP struct{}

// GenerateClient generates an HTTP/SOCKS client configuration
func (h *HTTP) GenerateClient(email, username, password string) models.InboundClient {
	return models.InboundClient{
		ID:     username,
		Email:  email,
		Enable: true,
	}
}

// GenerateShareLink generates an HTTP proxy link
func (h *HTTP) GenerateShareLink(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	// http://username:password@address:port
	if inbound.Settings.Password != "" {
		link := fmt.Sprintf("http://%s:%s@%s:%d",
			client.ID,
			inbound.Settings.Password,
			address,
			inbound.Port,
		)
		return link, nil
	}

	link := fmt.Sprintf("http://%s:%d", address, inbound.Port)
	return link, nil
}

// GenerateSOCKSLink generates a SOCKS5 proxy link
func (h *HTTP) GenerateSOCKSLink(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	// socks5://username:password@address:port
	if client.ID != "" && inbound.Settings.Password != "" {
		link := fmt.Sprintf("socks5://%s:%s@%s:%d",
			client.ID,
			inbound.Settings.Password,
			address,
			inbound.Port,
		)
		return link, nil
	}

	link := fmt.Sprintf("socks5://%s:%d", address, inbound.Port)
	return link, nil
}

// GenerateQRData generates QR code data
func (h *HTTP) GenerateQRData(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	return h.GenerateShareLink(inbound, client, address)
}

// ValidateSettings validates HTTP/SOCKS settings
func (h *HTTP) ValidateSettings(settings *models.InboundSettings) error {
	// HTTP/SOCKS can work without authentication
	return nil
}

// GetDefaultSettings returns default HTTP settings
func (h *HTTP) GetDefaultSettings() models.InboundSettings {
	return models.InboundSettings{
		Clients: []models.InboundClient{},
	}
}

// GeneratePACFile generates a PAC (Proxy Auto-Config) file
func (h *HTTP) GeneratePACFile(address string, port int) string {
	pac := fmt.Sprintf(`function FindProxyForURL(url, host) {
    // Direct access to local addresses
    if (isPlainHostName(host) ||
        shExpMatch(host, "*.local") ||
        isInNet(dnsResolve(host), "10.0.0.0", "255.0.0.0") ||
        isInNet(dnsResolve(host), "172.16.0.0",  "255.240.0.0") ||
        isInNet(dnsResolve(host), "192.168.0.0",  "255.255.0.0") ||
        isInNet(dnsResolve(host), "127.0.0.0", "255.255.255.0"))
        return "DIRECT";

    // All other traffic goes through proxy
    return "PROXY %s:%d";
}`, address, port)

	return pac
}

// GenerateCurlCommand generates a curl command for testing
func (h *HTTP) GenerateCurlCommand(inbound *models.Inbound, client *models.InboundClient, address string) string {
	if client.ID != "" && inbound.Settings.Password != "" {
		return fmt.Sprintf("curl -x http://%s:%s@%s:%d https://ipinfo.io",
			client.ID,
			inbound.Settings.Password,
			address,
			inbound.Port,
		)
	}

	return fmt.Sprintf("curl -x http://%s:%d https://ipinfo.io", address, inbound.Port)
}
