package protocols

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/v-ui/backend/internal/models"
	"golang.org/x/crypto/curve25519"
)

// WireGuard protocol handler
type WireGuard struct{}

// GenerateKeyPair generates a WireGuard private/public key pair
func (w *WireGuard) GenerateKeyPair() (privateKey, publicKey string, err error) {
	// Generate private key
	var privKey [32]byte
	if _, err := rand.Read(privKey[:]); err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Clamp the private key
	privKey[0] &= 248
	privKey[31] &= 127
	privKey[31] |= 64

	// Generate public key
	var pubKey [32]byte
	curve25519.ScalarBaseMult(&pubKey, &privKey)

	// Encode to base64
	privateKey = base64.StdEncoding.EncodeToString(privKey[:])
	publicKey = base64.StdEncoding.EncodeToString(pubKey[:])

	return privateKey, publicKey, nil
}

// GenerateClient generates a WireGuard client configuration
func (w *WireGuard) GenerateClient(email string) (models.InboundClient, string, error) {
	privateKey, publicKey, err := w.GenerateKeyPair()
	if err != nil {
		return models.InboundClient{}, "", err
	}

	client := models.InboundClient{
		ID:     publicKey, // WireGuard uses public key as ID
		Email:  email,
		Enable: true,
	}

	return client, privateKey, nil
}

// GenerateConfig generates WireGuard client config file
func (w *WireGuard) GenerateConfig(inbound *models.Inbound, client *models.InboundClient, privateKey, serverAddress string) (string, error) {
	// WireGuard config format
	config := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = 10.0.0.2/24
DNS = 1.1.1.1, 8.8.8.8

[Peer]
PublicKey = %s
Endpoint = %s:%d
AllowedIPs = 0.0.0.0/0
PersistentKeepalive = 25
`,
		privateKey,
		inbound.Settings.Password, // Server public key stored in password field
		serverAddress,
		inbound.Port,
	)

	return config, nil
}

// GenerateShareLink generates a WireGuard share link (not standard, custom format)
func (w *WireGuard) GenerateShareLink(inbound *models.Inbound, client *models.InboundClient, privateKey, address string) (string, error) {
	// WireGuard doesn't have a standard share link format
	// We'll create a custom one: wg://base64(config)
	config, err := w.GenerateConfig(inbound, client, privateKey, address)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(config))
	return "wg://" + encoded, nil
}

// ValidateSettings validates WireGuard settings
func (w *WireGuard) ValidateSettings(settings *models.InboundSettings) error {
	if settings.Password == "" {
		return fmt.Errorf("server public key (password field) is required")
	}

	// Validate base64 format
	if _, err := base64.StdEncoding.DecodeString(settings.Password); err != nil {
		return fmt.Errorf("invalid server public key format: %w", err)
	}

	return nil
}

// GetDefaultSettings returns default WireGuard settings
func (w *WireGuard) GetDefaultSettings() models.InboundSettings {
	return models.InboundSettings{
		Clients: []models.InboundClient{},
	}
}
