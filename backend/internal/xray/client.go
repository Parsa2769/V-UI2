package xray

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Client manages Xray core operations
type Client struct {
	configPath string
	logger     *zap.Logger
	cmd        *exec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewClient creates a new Xray client
func NewClient(configPath string, logger *zap.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		configPath: configPath,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Config represents Xray configuration
type Config struct {
	Log       LogConfig        `json:"log"`
	Inbounds  []InboundConfig  `json:"inbounds"`
	Outbounds []OutboundConfig `json:"outbounds"`
	Routing   RoutingConfig    `json:"routing"`
}

type LogConfig struct {
	Loglevel string `json:"loglevel"`
	Access   string `json:"access"`
	Error    string `json:"error"`
}

type InboundConfig struct {
	Port     int                    `json:"port"`
	Protocol string                 `json:"protocol"`
	Settings map[string]interface{} `json:"settings"`
	Tag      string                 `json:"tag"`
}

type OutboundConfig struct {
	Protocol string                 `json:"protocol"`
	Settings map[string]interface{} `json:"settings"`
	Tag      string                 `json:"tag"`
}

type RoutingConfig struct {
	Rules []RoutingRule `json:"rules"`
}

type RoutingRule struct {
	Type        string   `json:"type"`
	OutboundTag string   `json:"outboundTag"`
	Domain      []string `json:"domain,omitempty"`
}

// Start starts the Xray core process
func (c *Client) Start() error {
	if !c.IsConfigExists() {
		c.logger.Info("Xray config not found, creating default config")
		if err := c.CreateDefaultConfig(); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	c.cmd = exec.CommandContext(c.ctx, "xray", "run", "-c", c.configPath)
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start xray: %w", err)
	}

	c.logger.Info("Xray core started", zap.Int("pid", c.cmd.Process.Pid))

	// Monitor process in background
	go func() {
		if err := c.cmd.Wait(); err != nil {
			c.logger.Error("Xray process exited", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the Xray core process
func (c *Client) Stop() error {
	c.cancel()
	if c.cmd != nil && c.cmd.Process != nil {
		if err := c.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill xray process: %w", err)
		}
		c.logger.Info("Xray core stopped")
	}
	return nil
}

// Restart restarts the Xray core
func (c *Client) Restart() error {
	if err := c.Stop(); err != nil {
		c.logger.Warn("Error stopping xray", zap.Error(err))
	}
	time.Sleep(time.Second)
	return c.Start()
}

// IsRunning checks if Xray is running
func (c *Client) IsRunning() bool {
	if c.cmd == nil || c.cmd.Process == nil {
		return false
	}
	// Try to signal the process (0 = check if exists)
	return c.cmd.Process.Signal(os.Signal(nil)) == nil
}

// IsConfigExists checks if config file exists
func (c *Client) IsConfigExists() bool {
	_, err := os.Stat(c.configPath)
	return err == nil
}

// CreateDefaultConfig creates a default Xray configuration
func (c *Client) CreateDefaultConfig() error {
	config := Config{
		Log: LogConfig{
			Loglevel: "warning",
			Access:   "/var/log/xray/access.log",
			Error:    "/var/log/xray/error.log",
		},
		Inbounds: []InboundConfig{
			{
				Port:     10085,
				Protocol: "dokodemo-door",
				Settings: map[string]interface{}{
					"address": "0.0.0.0",
				},
				Tag: "api",
			},
		},
		Outbounds: []OutboundConfig{
			{
				Protocol: "freedom",
				Settings: map[string]interface{}{},
				Tag:      "direct",
			},
			{
				Protocol: "blackhole",
				Settings: map[string]interface{}{},
				Tag:      "blocked",
			},
		},
		Routing: RoutingConfig{
			Rules: []RoutingRule{
				{
					Type:        "field",
					OutboundTag: "api",
					Domain:      []string{"geosite:category-ads-all"},
				},
			},
		},
	}

	// Ensure directory exists
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	c.logger.Info("Default Xray config created", zap.String("path", c.configPath))
	return nil
}

// UpdateConfig updates the Xray configuration
func (c *Client) UpdateConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	c.logger.Info("Xray config updated")
	return c.Restart()
}

// GetConfig reads the current Xray configuration
func (c *Client) GetConfig() (*Config, error) {
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetVersion returns the Xray version
func (c *Client) GetVersion() (string, error) {
	cmd := exec.Command("xray", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get xray version: %w", err)
	}
	return string(output), nil
}

