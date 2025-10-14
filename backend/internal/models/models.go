package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

type Protocol string

const (
	ProtocolVMESS       Protocol = "vmess"
	ProtocolVLESS       Protocol = "vless"
	ProtocolTrojan      Protocol = "trojan"
	ProtocolShadowsocks Protocol = "shadowsocks"
	ProtocolWireGuard   Protocol = "wireguard"
)

type NodeStatus string

const (
	NodeStatusOnline  NodeStatus = "online"
	NodeStatusOffline NodeStatus = "offline"
	NodeStatusError   NodeStatus = "error"
)

// User represents a system user (admin, operator, etc.)
type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username        string     `gorm:"uniqueIndex;not null" json:"username"`
	Email           string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash    string     `gorm:"not null" json:"-"`
	Role            Role       `gorm:"not null;default:'viewer'" json:"role"`
	Enabled         bool       `gorm:"not null;default:true" json:"enabled"`
	TwoFactorSecret string     `json:"-"`
	TwoFactorEnabled bool      `gorm:"default:false" json:"two_factor_enabled"`
	LastLogin       *time.Time `json:"last_login"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Client represents an Xray client/end-user
type Client struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	UUID         string     `gorm:"uniqueIndex;not null" json:"uuid"`
	Protocol     Protocol   `gorm:"not null" json:"protocol"`
	Enable       bool       `gorm:"not null;default:true" json:"enable"`
	ExpiryTime   *time.Time `json:"expiry_time"`
	TrafficLimit int64      `gorm:"default:0" json:"traffic_limit"` // bytes, 0 = unlimited
	UploadBytes  int64      `gorm:"default:0" json:"upload_bytes"`
	DownloadBytes int64     `gorm:"default:0" json:"download_bytes"`
	IPLimit      int        `gorm:"default:0" json:"ip_limit"` // 0 = unlimited
	NodeID       *uuid.UUID `gorm:"type:uuid" json:"node_id"`
	Node         *Node      `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	CreatedBy    uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Node represents an Xray server node
type Node struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name       string     `gorm:"uniqueIndex;not null" json:"name"`
	Address    string     `gorm:"not null" json:"address"`
	Port       int        `gorm:"not null" json:"port"`
	Status     NodeStatus `gorm:"not null;default:'offline'" json:"status"`
	LastSeen   *time.Time `json:"last_seen"`
	UploadBytes int64     `gorm:"default:0" json:"upload_bytes"`
	DownloadBytes int64   `gorm:"default:0" json:"download_bytes"`
	Metadata   string     `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TrafficLog represents traffic usage records
type TrafficLog struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ClientID      uuid.UUID `gorm:"type:uuid;not null;index" json:"client_id"`
	Client        *Client   `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	NodeID        *uuid.UUID `gorm:"type:uuid;index" json:"node_id"`
	Node          *Node     `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	UploadBytes   int64     `gorm:"not null" json:"upload_bytes"`
	DownloadBytes int64     `gorm:"not null" json:"download_bytes"`
	Timestamp     time.Time `gorm:"not null;index" json:"timestamp"`
}

// AuditLog represents security audit logs
type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ActorID   uuid.UUID `gorm:"type:uuid;not null;index" json:"actor_id"`
	Actor     *User     `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
	Action    string    `gorm:"not null;index" json:"action"`
	Target    string    `gorm:"not null" json:"target"`
	Detail    string    `gorm:"type:text" json:"detail"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `gorm:"not null;index" json:"created_at"`
}

// RefreshToken represents JWT refresh tokens
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

// ConfigTemplate represents pre-defined configuration templates
type ConfigTemplate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Protocol    Protocol  `gorm:"not null" json:"protocol"`
	Description string    `gorm:"type:text" json:"description"`
	Template    string    `gorm:"type:text;not null" json:"template"`
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName methods for custom table names
func (User) TableName() string { return "users" }
func (Client) TableName() string { return "clients" }
func (Node) TableName() string { return "nodes" }
func (TrafficLog) TableName() string { return "traffic_logs" }
func (AuditLog) TableName() string { return "audit_logs" }
func (RefreshToken) TableName() string { return "refresh_tokens" }
func (ConfigTemplate) TableName() string { return "config_templates" }

