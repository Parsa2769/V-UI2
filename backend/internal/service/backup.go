package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

type BackupService struct {
	db         *db.Database
	log        *zap.Logger
	backupPath string
}

func NewBackupService(database *db.Database, log *zap.Logger, backupPath string) *BackupService {
	return &BackupService{
		db:         database,
		log:        log,
		backupPath: backupPath,
	}
}

// BackupData represents the complete backup structure
type BackupData struct {
	Version    string                `json:"version"`
	Timestamp  time.Time             `json:"timestamp"`
	Users      []models.User         `json:"users"`
	Clients    []models.Client       `json:"clients"`
	Nodes      []models.Node         `json:"nodes"`
	Inbounds   []models.Inbound      `json:"inbounds"`
	TrafficLogs []models.TrafficLog  `json:"traffic_logs"`
	AuditLogs  []models.AuditLog     `json:"audit_logs"`
	Templates  []models.ConfigTemplate `json:"templates"`
}

// ExportDatabase exports the entire database to JSON
func (s *BackupService) ExportDatabase() (*BackupData, error) {
	s.log.Info("Exporting database...")

	backup := &BackupData{
		Version:   "2.0.0",
		Timestamp: time.Now(),
	}

	// Export users
	if err := s.db.Find(&backup.Users).Error; err != nil {
		return nil, fmt.Errorf("failed to export users: %w", err)
	}

	// Export clients
	if err := s.db.Find(&backup.Clients).Error; err != nil {
		return nil, fmt.Errorf("failed to export clients: %w", err)
	}

	// Export nodes
	if err := s.db.Find(&backup.Nodes).Error; err != nil {
		return nil, fmt.Errorf("failed to export nodes: %w", err)
	}

	// Export inbounds
	if err := s.db.Find(&backup.Inbounds).Error; err != nil {
		return nil, fmt.Errorf("failed to export inbounds: %w", err)
	}

	// Export traffic logs (last 30 days only)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if err := s.db.Where("created_at > ?", thirtyDaysAgo).Find(&backup.TrafficLogs).Error; err != nil {
		return nil, fmt.Errorf("failed to export traffic logs: %w", err)
	}

	// Export audit logs (last 90 days only)
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	if err := s.db.Where("timestamp > ?", ninetyDaysAgo).Find(&backup.AuditLogs).Error; err != nil {
		return nil, fmt.Errorf("failed to export audit logs: %w", err)
	}

	// Export templates
	if err := s.db.Find(&backup.Templates).Error; err != nil {
		return nil, fmt.Errorf("failed to export templates: %w", err)
	}

	s.log.Info("Database exported successfully",
		zap.Int("users", len(backup.Users)),
		zap.Int("clients", len(backup.Clients)),
		zap.Int("nodes", len(backup.Nodes)),
		zap.Int("inbounds", len(backup.Inbounds)))

	return backup, nil
}

// CreateBackup creates a full backup file
func (s *BackupService) CreateBackup() (string, error) {
	s.log.Info("Creating backup...")

	// Export data
	data, err := s.ExportDatabase()
	if err != nil {
		return "", err
	}

	// Create backup directory
	if err := os.MkdirAll(s.backupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("v-ui-backup-%s.zip", time.Now().Format("2006-01-02-15-04-05"))
	backupFile := filepath.Join(s.backupPath, filename)

	// Create zip file
	zipFile, err := os.Create(backupFile)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Write database JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	dbFile, err := zipWriter.Create("database.json")
	if err != nil {
		return "", fmt.Errorf("failed to create database.json in zip: %w", err)
	}

	if _, err := dbFile.Write(jsonData); err != nil {
		return "", fmt.Errorf("failed to write database.json: %w", err)
	}

	// Add metadata
	metadata := map[string]interface{}{
		"version":     "2.0.0",
		"timestamp":   time.Now(),
		"record_count": map[string]int{
			"users":      len(data.Users),
			"clients":    len(data.Clients),
			"nodes":      len(data.Nodes),
			"inbounds":   len(data.Inbounds),
			"traffic_logs": len(data.TrafficLogs),
			"audit_logs":   len(data.AuditLogs),
		},
	}

	metadataJSON, _ := json.MarshalIndent(metadata, "", "  ")
	metaFile, _ := zipWriter.Create("metadata.json")
	metaFile.Write(metadataJSON)

	s.log.Info("Backup created successfully", zap.String("file", backupFile))
	return backupFile, nil
}

// ImportDatabase imports data from backup
func (s *BackupService) ImportDatabase(data *BackupData, options ImportOptions) error {
	s.log.Info("Importing database...", zap.Any("options", options))

	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Import users
	if options.ImportUsers && len(data.Users) > 0 {
		s.log.Info("Importing users", zap.Int("count", len(data.Users)))
		for _, user := range data.Users {
			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to import user %s: %w", user.Username, err)
			}
		}
	}

	// Import clients
	if options.ImportClients && len(data.Clients) > 0 {
		s.log.Info("Importing clients", zap.Int("count", len(data.Clients)))
		for _, client := range data.Clients {
			if err := tx.Create(&client).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to import client %s: %w", client.Email, err)
			}
		}
	}

	// Import nodes
	if options.ImportNodes && len(data.Nodes) > 0 {
		s.log.Info("Importing nodes", zap.Int("count", len(data.Nodes)))
		for _, node := range data.Nodes {
			if err := tx.Create(&node).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to import node %s: %w", node.Name, err)
			}
		}
	}

	// Import inbounds
	if options.ImportInbounds && len(data.Inbounds) > 0 {
		s.log.Info("Importing inbounds", zap.Int("count", len(data.Inbounds)))
		for _, inbound := range data.Inbounds {
			if err := tx.Create(&inbound).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to import inbound %s: %w", inbound.Tag, err)
			}
		}
	}

	// Import templates
	if options.ImportTemplates && len(data.Templates) > 0 {
		s.log.Info("Importing templates", zap.Int("count", len(data.Templates)))
		for _, template := range data.Templates {
			if err := tx.Create(&template).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to import template %s: %w", template.Name, err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.log.Info("Database imported successfully")
	return nil
}

// RestoreFromFile restores database from a backup file
func (s *BackupService) RestoreFromFile(backupFile string, options ImportOptions) error {
	s.log.Info("Restoring from backup file", zap.String("file", backupFile))

	// Open zip file
	reader, err := zip.OpenReader(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer reader.Close()

	// Find database.json
	var jsonFile *zip.File
	for _, file := range reader.File {
		if file.Name == "database.json" {
			jsonFile = file
			break
		}
	}

	if jsonFile == nil {
		return fmt.Errorf("database.json not found in backup")
	}

	// Read database.json
	rc, err := jsonFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open database.json: %w", err)
	}
	defer rc.Close()

	jsonData, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("failed to read database.json: %w", err)
	}

	// Parse JSON
	var data BackupData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to parse backup data: %w", err)
	}

	// Import data
	return s.ImportDatabase(&data, options)
}

// ImportOptions specifies what to import
type ImportOptions struct {
	ImportUsers     bool
	ImportClients   bool
	ImportNodes     bool
	ImportInbounds  bool
	ImportTemplates bool
	Overwrite       bool
}

// ListBackups lists all available backups
func (s *BackupService) ListBackups() ([]BackupInfo, error) {
	files, err := os.ReadDir(s.backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupInfo
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".zip" {
			info, _ := file.Info()
			backups = append(backups, BackupInfo{
				Filename:  file.Name(),
				Size:      info.Size(),
				CreatedAt: info.ModTime(),
			})
		}
	}

	return backups, nil
}

// BackupInfo represents backup file information
type BackupInfo struct {
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// DeleteBackup deletes a backup file
func (s *BackupService) DeleteBackup(filename string) error {
	backupFile := filepath.Join(s.backupPath, filename)
	if err := os.Remove(backupFile); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}
	s.log.Info("Backup deleted", zap.String("file", filename))
	return nil
}
