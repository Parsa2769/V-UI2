package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Inbound represents an inbound configuration
type Inbound struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Tag       string         `gorm:"uniqueIndex;not null" json:"tag"`
	Protocol  string         `gorm:"not null" json:"protocol"` // vmess, vless, trojan, shadowsocks, etc.
	Port      int            `gorm:"not null" json:"port"`
	Listen    string         `gorm:"default:0.0.0.0" json:"listen"`
	Enable    bool           `gorm:"default:true" json:"enable"`
	Settings  InboundSettings `gorm:"type:jsonb" json:"settings"`
	StreamSettings StreamSettings `gorm:"type:jsonb" json:"stream_settings"`
	Sniffing  SniffingSettings `gorm:"type:jsonb" json:"sniffing"`
	Allocate  AllocateSettings `gorm:"type:jsonb" json:"allocate"`
	Remark    string         `json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// InboundSettings holds protocol-specific settings
type InboundSettings struct {
	Clients      []InboundClient `json:"clients,omitempty"`
	Decryption   string          `json:"decryption,omitempty"`
	Fallbacks    []Fallback      `json:"fallbacks,omitempty"`
	Network      string          `json:"network,omitempty"`
	Password     string          `json:"password,omitempty"`
	Method       string          `json:"method,omitempty"`
	Level        int             `json:"level,omitempty"`
}

// InboundClient represents a client in inbound settings
type InboundClient struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Flow       string `json:"flow,omitempty"`
	LimitIP    int    `json:"limitIp,omitempty"`
	TotalGB    int64  `json:"totalGB,omitempty"`
	ExpiryTime int64  `json:"expiryTime,omitempty"`
	Enable     bool   `json:"enable"`
	TgID       string `json:"tgId,omitempty"`
	SubID      string `json:"subId,omitempty"`
}

// StreamSettings for network transport
type StreamSettings struct {
	Network      string            `json:"network"` // tcp, ws, grpc, http, quic, kcp
	Security     string            `json:"security"` // none, tls, reality
	TLSSettings  *TLSSettings      `json:"tlsSettings,omitempty"`
	TCPSettings  *TCPSettings      `json:"tcpSettings,omitempty"`
	WSSettings   *WSSettings       `json:"wsSettings,omitempty"`
	HTTPSettings *HTTPSettings     `json:"httpSettings,omitempty"`
	GRPCSettings *GRPCSettings     `json:"grpcSettings,omitempty"`
	QUICSettings *QUICSettings     `json:"quicSettings,omitempty"`
	KCPSettings  *KCPSettings      `json:"kcpSettings,omitempty"`
	RealitySettings *RealitySettings `json:"realitySettings,omitempty"`
}

// TLSSettings for TLS encryption
type TLSSettings struct {
	ServerName    string   `json:"serverName,omitempty"`
	Certificates  []Certificate `json:"certificates,omitempty"`
	ALPN          []string `json:"alpn,omitempty"`
	AllowInsecure bool     `json:"allowInsecure,omitempty"`
}

// RealitySettings for XTLS Reality
type RealitySettings struct {
	Show          bool     `json:"show"`
	Dest          string   `json:"dest"`
	Xver          int      `json:"xver"`
	ServerNames   []string `json:"serverNames"`
	PrivateKey    string   `json:"privateKey"`
	ShortIds      []string `json:"shortIds"`
	Settings      RealityInnerSettings `json:"settings"`
}

// RealityInnerSettings inner settings for Reality
type RealityInnerSettings struct {
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
	ServerName  string `json:"serverName"`
	SpiderX     string `json:"spiderX"`
}

// Certificate holds TLS certificate info
type Certificate struct {
	CertificateFile string `json:"certificateFile"`
	KeyFile         string `json:"keyFile"`
	Certificate     []string `json:"certificate,omitempty"`
	Key             []string `json:"key,omitempty"`
}

// TCPSettings for TCP transport
type TCPSettings struct {
	Header HTTPHeader `json:"header,omitempty"`
}

// WSSettings for WebSocket transport
type WSSettings struct {
	Path    string     `json:"path"`
	Headers HTTPHeader `json:"headers,omitempty"`
}

// HTTPSettings for HTTP/2 transport
type HTTPSettings struct {
	Host []string `json:"host,omitempty"`
	Path string   `json:"path,omitempty"`
}

// GRPCSettings for gRPC transport
type GRPCSettings struct {
	ServiceName string `json:"serviceName"`
	MultiMode   bool   `json:"multiMode,omitempty"`
}

// QUICSettings for QUIC transport
type QUICSettings struct {
	Security string     `json:"security"`
	Key      string     `json:"key,omitempty"`
	Header   HTTPHeader `json:"header,omitempty"`
}

// KCPSettings for mKCP transport
type KCPSettings struct {
	MTU              int        `json:"mtu,omitempty"`
	TTI              int        `json:"tti,omitempty"`
	UplinkCapacity   int        `json:"uplinkCapacity,omitempty"`
	DownlinkCapacity int        `json:"downlinkCapacity,omitempty"`
	Congestion       bool       `json:"congestion,omitempty"`
	ReadBufferSize   int        `json:"readBufferSize,omitempty"`
	WriteBufferSize  int        `json:"writeBufferSize,omitempty"`
	Header           HTTPHeader `json:"header,omitempty"`
}

// HTTPHeader for various transport headers
type HTTPHeader struct {
	Type    string              `json:"type,omitempty"`
	Request *HTTPRequestHeader  `json:"request,omitempty"`
	Response *HTTPResponseHeader `json:"response,omitempty"`
}

// HTTPRequestHeader for HTTP request headers
type HTTPRequestHeader struct {
	Version string              `json:"version,omitempty"`
	Method  string              `json:"method,omitempty"`
	Path    []string            `json:"path,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// HTTPResponseHeader for HTTP response headers
type HTTPResponseHeader struct {
	Version string              `json:"version,omitempty"`
	Status  string              `json:"status,omitempty"`
	Reason  string              `json:"reason,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// SniffingSettings for traffic sniffing
type SniffingSettings struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

// AllocateSettings for port allocation
type AllocateSettings struct {
	Strategy    string `json:"strategy,omitempty"`
	Refresh     int    `json:"refresh,omitempty"`
	Concurrency int    `json:"concurrency,omitempty"`
}

// Fallback for trojan fallback
type Fallback struct {
	Name string `json:"name,omitempty"`
	Alpn string `json:"alpn,omitempty"`
	Path string `json:"path,omitempty"`
	Dest string `json:"dest"`
	Xver int    `json:"xver,omitempty"`
}

// Scan implementations for GORM
func (s *InboundSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s InboundSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StreamSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s StreamSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SniffingSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s SniffingSettings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (a *AllocateSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

func (a AllocateSettings) Value() (driver.Value, error) {
	return json.Marshal(a)
}

