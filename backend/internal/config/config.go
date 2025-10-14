package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Auth       AuthConfig       `yaml:"auth"`
	Log        LogConfig        `yaml:"log"`
	Xray       XrayConfig       `yaml:"xray"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	RateLimit  RateLimitConfig  `yaml:"rate_limit"`
	CORS       CORSConfig       `yaml:"cors"`
	Session    SessionConfig    `yaml:"session"`
	Backup     BackupConfig     `yaml:"backup"`
	Notification NotificationConfig `yaml:"notification"`
	Defaults   DefaultsConfig   `yaml:"defaults"`
	Telegram   TelegramConfig   `yaml:"telegram"`
	Certificate CertConfig      `yaml:"certificate"`
}

type ServerConfig struct {
	Host string
	Port int
	Mode string
}

type DatabaseConfig struct {
	Type         string
	Host         string
	Port         int
	Name         string
	User         string
	Password     string
	SSLMode      string
	Path         string
	MaxOpenConns int
	MaxIdleConns int
}

type AuthConfig struct {
	JWTSecret        string
	JWTRefreshSecret string
	AccessExpiry     time.Duration
	RefreshExpiry    time.Duration
	Enable2FA        bool
	BcryptCost       int
}

type XrayConfig struct {
	ConfigPath string
	APIHost    string
	APIPort    int
	AssetsPath string
}

type MonitoringConfig struct {
	EnableMetrics   bool
	EnableWebSocket bool
	MetricsPath     string
	WebSocketPath   string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type LogConfig struct {
	Level  string
	Format string
	Output string
}

type RateLimitConfig struct {
	Enabled  bool
	Requests int
	Window   time.Duration
}

func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read from environment
	v.AutomaticEnv()

	// Read config file if exists
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate config
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Server
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "production")

	// Database
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.name", "x3ui")
	v.SetDefault("database.user", "x3ui")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.path", "./data/x3ui.db")
	v.SetDefault("database.maxopenconns", 25)
	v.SetDefault("database.maxidleconns", 5)

	// Auth
	v.SetDefault("auth.accessexpiry", "24h")
	v.SetDefault("auth.refreshexpiry", "168h")
	v.SetDefault("auth.enable2fa", true)
	v.SetDefault("auth.bcryptcost", 12)

	// Xray
	v.SetDefault("xray.configpath", "/etc/xray/config.json")
	v.SetDefault("xray.apihost", "localhost")
	v.SetDefault("xray.apiport", 10085)
	v.SetDefault("xray.assetspath", "/usr/local/share/xray")

	// Monitoring
	v.SetDefault("monitoring.enablemetrics", true)
	v.SetDefault("monitoring.enablewebsocket", true)
	v.SetDefault("monitoring.metricspath", "/metrics")
	v.SetDefault("monitoring.websocketpath", "/ws/monitor")

	// CORS
	v.SetDefault("cors.allowedorigins", []string{"http://localhost:3000"})
	v.SetDefault("cors.allowedmethods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowedheaders", []string{"Content-Type", "Authorization"})
	v.SetDefault("cors.allowcredentials", true)

	// Logging
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")

	// Rate Limit
	v.SetDefault("ratelimit.enabled", true)
	v.SetDefault("ratelimit.requests", 100)
	v.SetDefault("ratelimit.window", "1m")
}

func validateConfig(cfg *Config) error {
	if cfg.Auth.JWTSecret == "" || cfg.Auth.JWTSecret == "changeme_32_char_random_secret_key_here_minimum" {
		return fmt.Errorf("JWT_SECRET must be set to a secure random value")
	}

	if len(cfg.Auth.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if cfg.Auth.JWTRefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET must be set")
	}

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	return nil
}

