package service

import (
	"github.com/v-ui/backend/internal/auth"
	"github.com/v-ui/backend/internal/config"
	"github.com/v-ui/backend/internal/db"
	"go.uber.org/zap"
)

type Services struct {
	Auth         *AuthService
	User         *UserService
	Client       *ClientService
	Node         *NodeService
	Traffic      *TrafficService
	Audit        *AuditService
	Inbound      *InboundService
	Subscription *SubscriptionService
	Template     *TemplateService
	Backup       *BackupService
}

func NewServices(database *db.Database, cfg *config.Config, log *zap.Logger) *Services {
	jwtManager := auth.NewJWTManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTRefreshSecret,
		cfg.Auth.AccessExpiry,
		cfg.Auth.RefreshExpiry,
	)

	return &Services{
		Auth:         NewAuthService(database, jwtManager, cfg, log),
		User:         NewUserService(database, cfg, log),
		Client:       NewClientService(database, log),
		Node:         NewNodeService(database, log),
		Traffic:      NewTrafficService(database, log),
		Audit:        NewAuditService(database, log),
		Inbound:      NewInboundService(database, log),
		Subscription: NewSubscriptionService(database, log),
		Template:     NewTemplateService(database, log),
		Backup:       NewBackupService(database, log, cfg.Backup.Path),
	}
}

