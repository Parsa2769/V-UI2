package protocols

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/v-ui/backend/internal/models"
)

// Shadowsocks protocol handler
type Shadowsocks struct{}

// GenerateClient generates a Shadowsocks client configuration
func (s *Shadowsocks) GenerateClient(email, password, method string) models.InboundClient {
	return models.InboundClient{
		ID:     password,
		Email:  email,
		Enable: true,
	}
}

// GenerateShareLink generates an ss:// share link
func (s *Shadowsocks) GenerateShareLink(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	// ss://base64(method:password)@address:port#remark
	
	method := inbound.Settings.Method
	if method == "" {
		method = "aes-256-gcm" // default
	}

	// Encode method:password
	userInfo := fmt.Sprintf("%s:%s", method, client.ID)
	encoded := base64.URLEncoding.EncodeToString([]byte(userInfo))

	remark := url.QueryEscape(client.Email)
	link := fmt.Sprintf("ss://%s@%s:%d#%s",
		encoded,
		address,
		inbound.Port,
		remark,
	)

	return link, nil
}

// GenerateQRData generates QR code data
func (s *Shadowsocks) GenerateQRData(inbound *models.Inbound, client *models.InboundClient, address string) (string, error) {
	return s.GenerateShareLink(inbound, client, address)
}

// ValidateSettings validates Shadowsocks settings
func (s *Shadowsocks) ValidateSettings(settings *models.InboundSettings) error {
	validMethods := []string{
		"aes-128-gcm",
		"aes-256-gcm",
		"chacha20-poly1305",
		"chacha20-ietf-poly1305",
		"xchacha20-poly1305",
		"xchacha20-ietf-poly1305",
		"2022-blake3-aes-128-gcm",
		"2022-blake3-aes-256-gcm",
		"2022-blake3-chacha20-poly1305",
	}

	if settings.Method == "" {
		return fmt.Errorf("encryption method is required")
	}

	valid := false
	for _, m := range validMethods {
		if settings.Method == m {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid encryption method: %s", settings.Method)
	}

	if settings.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// GetDefaultSettings returns default Shadowsocks settings
func (s *Shadowsocks) GetDefaultSettings() models.InboundSettings {
	return models.InboundSettings{
		Method:   "aes-256-gcm",
		Password: "",
		Network:  "tcp,udp",
	}
}
