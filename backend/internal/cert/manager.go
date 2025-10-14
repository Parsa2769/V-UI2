package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

// Manager handles SSL/TLS certificate operations
type Manager struct {
	certDir    string
	acmeEmail  string
	logger     *zap.Logger
	autocertMgr *autocert.Manager
}

// NewManager creates a new certificate manager
func NewManager(certDir, acmeEmail string, logger *zap.Logger) *Manager {
	return &Manager{
		certDir:   certDir,
		acmeEmail: acmeEmail,
		logger:    logger,
	}
}

// GenerateSelfSigned generates a self-signed certificate
func (m *Manager) GenerateSelfSigned(domain string, validDays int) (certPath, keyPath string, err error) {
	m.logger.Info("Generating self-signed certificate", zap.String("domain", domain))

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(validDays) * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"V-UI"},
			CommonName:   domain,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{domain},
	}

	// Create certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %w", err)
	}

	// Ensure cert directory exists
	if err := os.MkdirAll(m.certDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create cert directory: %w", err)
	}

	// Save certificate
	certPath = filepath.Join(m.certDir, domain+".crt")
	certOut, err := os.Create(certPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cert file: %w", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return "", "", fmt.Errorf("failed to write cert: %w", err)
	}

	// Save private key
	keyPath = filepath.Join(m.certDir, domain+".key")
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyOut.Close()

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return "", "", fmt.Errorf("failed to write key: %w", err)
	}

	m.logger.Info("Self-signed certificate generated",
		zap.String("domain", domain),
		zap.String("cert", certPath),
		zap.String("key", keyPath))

	return certPath, keyPath, nil
}

// RequestLetsEncrypt requests a certificate from Let's Encrypt
func (m *Manager) RequestLetsEncrypt(domains []string) error {
	m.logger.Info("Requesting Let's Encrypt certificate", zap.Strings("domains", domains))

	// Initialize autocert manager
	m.autocertMgr = &autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Email:       m.acmeEmail,
		HostPolicy:  autocert.HostWhitelist(domains...),
		Cache:       autocert.DirCache(m.certDir),
	}

	// This will be used by the HTTP/HTTPS server
	m.logger.Info("Let's Encrypt manager initialized")
	return nil
}

// GetAutocertManager returns the autocert manager for use with HTTP server
func (m *Manager) GetAutocertManager() *autocert.Manager {
	return m.autocertMgr
}

// RenewCertificate checks and renews certificate if needed
func (m *Manager) RenewCertificate(domain string) error {
	m.logger.Info("Checking certificate renewal", zap.String("domain", domain))

	certPath := filepath.Join(m.certDir, domain+".crt")

	// Load certificate
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %w", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Check if certificate expires in less than 30 days
	daysUntilExpiry := time.Until(cert.NotAfter).Hours() / 24
	if daysUntilExpiry > 30 {
		m.logger.Info("Certificate still valid",
			zap.String("domain", domain),
			zap.Float64("days_until_expiry", daysUntilExpiry))
		return nil
	}

	m.logger.Warn("Certificate needs renewal",
		zap.String("domain", domain),
		zap.Float64("days_until_expiry", daysUntilExpiry))

	// Renew using Let's Encrypt
	if m.autocertMgr != nil {
		// The autocert manager handles renewal automatically
		return nil
	}

	// Otherwise, generate a new self-signed certificate
	_, _, err = m.GenerateSelfSigned(domain, 365)
	return err
}

// GetCertificateInfo returns information about a certificate
func (m *Manager) GetCertificateInfo(domain string) (map[string]interface{}, error) {
	certPath := filepath.Join(m.certDir, domain+".crt")

	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	info := map[string]interface{}{
		"subject":     cert.Subject.String(),
		"issuer":      cert.Issuer.String(),
		"not_before":  cert.NotBefore,
		"not_after":   cert.NotAfter,
		"dns_names":   cert.DNSNames,
		"is_ca":       cert.IsCA,
		"serial":      cert.SerialNumber.String(),
	}

	return info, nil
}

// ListCertificates lists all certificates in the cert directory
func (m *Manager) ListCertificates() ([]string, error) {
	files, err := os.ReadDir(m.certDir)
	if err != nil {
		return nil, err
	}

	var certs []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".crt" {
			certs = append(certs, file.Name())
		}
	}

	return certs, nil
}

// DeleteCertificate deletes a certificate and its key
func (m *Manager) DeleteCertificate(domain string) error {
	certPath := filepath.Join(m.certDir, domain+".crt")
	keyPath := filepath.Join(m.certDir, domain+".key")

	if err := os.Remove(certPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete certificate: %w", err)
	}

	if err := os.Remove(keyPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	m.logger.Info("Certificate deleted", zap.String("domain", domain))
	return nil
}
