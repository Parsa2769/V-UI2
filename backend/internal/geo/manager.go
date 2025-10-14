package geo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

const (
	// Official Xray geo files
	GeoIPURL   = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat"
	GeoSiteURL = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"

	// Iran-specific rules
	GeoIPIranURL   = "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat"
	GeoSiteIranURL = "https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat"

	// Russia-specific rules
	GeoIPRussiaURL   = "https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat"
	GeoSiteRussiaURL = "https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat"
)

// Manager handles geo files download and updates
type Manager struct {
	assetsPath string
	logger     *zap.Logger
	client     *http.Client
}

// NewManager creates a new geo files manager
func NewManager(assetsPath string, logger *zap.Logger) *Manager {
	return &Manager{
		assetsPath: assetsPath,
		logger:     logger,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// DownloadFile downloads a file from URL to destination
func (m *Manager) DownloadFile(url, destination string) error {
	m.logger.Info("Downloading file", zap.String("url", url), zap.String("dest", destination))

	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download file
	resp, err := m.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create temp file
	tempFile := destination + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy content
	if _, err := io.Copy(out, resp.Body); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Rename temp file to destination
	if err := os.Rename(tempFile, destination); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename file: %w", err)
	}

	m.logger.Info("File downloaded successfully", zap.String("dest", destination))
	return nil
}

// UpdateGeoIP downloads the latest geoip.dat
func (m *Manager) UpdateGeoIP() error {
	destination := filepath.Join(m.assetsPath, "geoip.dat")
	return m.DownloadFile(GeoIPURL, destination)
}

// UpdateGeoSite downloads the latest geosite.dat
func (m *Manager) UpdateGeoSite() error {
	destination := filepath.Join(m.assetsPath, "geosite.dat")
	return m.DownloadFile(GeoSiteURL, destination)
}

// UpdateIranRules downloads Iran-specific routing rules
func (m *Manager) UpdateIranRules() error {
	m.logger.Info("Updating Iran routing rules...")

	if err := m.DownloadFile(GeoIPIranURL, filepath.Join(m.assetsPath, "geoip-iran.dat")); err != nil {
		return fmt.Errorf("failed to download geoip-iran: %w", err)
	}

	if err := m.DownloadFile(GeoSiteIranURL, filepath.Join(m.assetsPath, "geosite-iran.dat")); err != nil {
		return fmt.Errorf("failed to download geosite-iran: %w", err)
	}

	m.logger.Info("Iran rules updated successfully")
	return nil
}

// UpdateRussiaRules downloads Russia-specific routing rules
func (m *Manager) UpdateRussiaRules() error {
	m.logger.Info("Updating Russia routing rules...")

	if err := m.DownloadFile(GeoIPRussiaURL, filepath.Join(m.assetsPath, "geoip-russia.dat")); err != nil {
		return fmt.Errorf("failed to download geoip-russia: %w", err)
	}

	if err := m.DownloadFile(GeoSiteRussiaURL, filepath.Join(m.assetsPath, "geosite-russia.dat")); err != nil {
		return fmt.Errorf("failed to download geosite-russia: %w", err)
	}

	m.logger.Info("Russia rules updated successfully")
	return nil
}

// UpdateAll downloads all geo files
func (m *Manager) UpdateAll() error {
	m.logger.Info("Updating all geo files...")

	if err := m.UpdateGeoIP(); err != nil {
		m.logger.Error("Failed to update GeoIP", zap.Error(err))
		return err
	}

	if err := m.UpdateGeoSite(); err != nil {
		m.logger.Error("Failed to update GeoSite", zap.Error(err))
		return err
	}

	// Optional: Update region-specific rules
	// Uncomment if needed
	// m.UpdateIranRules()
	// m.UpdateRussiaRules()

	m.logger.Info("All geo files updated successfully")
	return nil
}

// CheckForUpdates checks if geo files need updating
func (m *Manager) CheckForUpdates() (bool, error) {
	geoipPath := filepath.Join(m.assetsPath, "geoip.dat")
	geositePath := filepath.Join(m.assetsPath, "geosite.dat")

	// Check if files exist
	geoipInfo, err := os.Stat(geoipPath)
	if err != nil {
		return true, nil // File doesn't exist, need update
	}

	geositeInfo, err := os.Stat(geositePath)
	if err != nil {
		return true, nil // File doesn't exist, need update
	}

	// Check if files are older than 7 days
	weekAgo := time.Now().AddDate(0, 0, -7)
	if geoipInfo.ModTime().Before(weekAgo) || geositeInfo.ModTime().Before(weekAgo) {
		return true, nil
	}

	return false, nil
}

// AutoUpdate runs automatic updates if needed
func (m *Manager) AutoUpdate() error {
	needsUpdate, err := m.CheckForUpdates()
	if err != nil {
		return err
	}

	if needsUpdate {
		m.logger.Info("Geo files need updating, starting auto-update...")
		return m.UpdateAll()
	}

	m.logger.Info("Geo files are up to date")
	return nil
}

// StartAutoUpdateScheduler starts a background scheduler for auto-updates
func (m *Manager) StartAutoUpdateScheduler(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			m.logger.Info("Running scheduled geo files update...")
			if err := m.AutoUpdate(); err != nil {
				m.logger.Error("Scheduled update failed", zap.Error(err))
			}
		}
	}()
}

// GetFileInfo returns information about geo files
func (m *Manager) GetFileInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	files := []string{"geoip.dat", "geosite.dat", "geoip-iran.dat", "geosite-iran.dat"}

	for _, file := range files {
		path := filepath.Join(m.assetsPath, file)
		stat, err := os.Stat(path)
		if err != nil {
			info[file] = map[string]interface{}{
				"exists": false,
			}
			continue
		}

		info[file] = map[string]interface{}{
			"exists":   true,
			"size":     stat.Size(),
			"mod_time": stat.ModTime(),
		}
	}

	return info, nil
}
