package api

import (
	"github.com/gin-gonic/gin"
	"github.com/v-ui/backend/internal/auth"
	"github.com/v-ui/backend/internal/config"
	"github.com/v-ui/backend/internal/middleware"
	"github.com/v-ui/backend/internal/service"
	"github.com/v-ui/backend/internal/websocket"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	engine   *gin.Engine
	services *service.Services
	cfg      *config.Config
	log      *zap.Logger
	wsHub    *websocket.Hub
}

func NewRouter(engine *gin.Engine, services *service.Services, cfg *config.Config, log *zap.Logger, wsHub *websocket.Hub) *Router {
	return &Router{
		engine:   engine,
		services: services,
		cfg:      cfg,
		log:      log,
		wsHub:    wsHub,
	}
}

func (r *Router) Setup() {
	// Swagger documentation
	r.engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Metrics endpoint
	if r.cfg.Monitoring.EnableMetrics {
		r.engine.GET(r.cfg.Monitoring.MetricsPath, gin.WrapH(promhttp.Handler()))
	}

	// WebSocket endpoint
	if r.cfg.Monitoring.EnableWebSocket && r.wsHub != nil {
		r.engine.GET(r.cfg.Monitoring.WebSocketPath, func(c *gin.Context) {
			websocket.ServeWs(r.wsHub, c.Writer, c.Request)
		})
	}

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		// Public auth endpoints
		authHandler := NewAuthHandler(r.services.Auth, r.services.Audit, r.log)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/2fa/verify", authHandler.Verify2FA)
		}

		// Create JWT manager for middleware
		jwtManager := auth.NewJWTManager(
			r.cfg.Auth.JWTSecret,
			r.cfg.Auth.JWTRefreshSecret,
			r.cfg.Auth.AccessExpiry,
			r.cfg.Auth.RefreshExpiry,
		)

		// Protected endpoints
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired(jwtManager))
		{
			// User management
			userHandler := NewUserHandler(r.services.User, r.services.Audit, r.log)
			users := protected.Group("/users")
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.GetByID)
				users.POST("", middleware.RequireRole("admin"), userHandler.Create)
				users.PUT("/:id", middleware.RequireRole("admin"), userHandler.Update)
				users.DELETE("/:id", middleware.RequireRole("admin"), userHandler.Delete)
			}

			// Client management
			clientHandler := NewClientHandler(r.services.Client, r.services.Audit, r.log)
			clients := protected.Group("/clients")
			{
				clients.GET("", clientHandler.List)
				clients.GET("/:id", clientHandler.GetByID)
				clients.POST("", middleware.RequireRole("admin", "operator"), clientHandler.Create)
				clients.PUT("/:id", middleware.RequireRole("admin", "operator"), clientHandler.Update)
				clients.DELETE("/:id", middleware.RequireRole("admin", "operator"), clientHandler.Delete)
				clients.GET("/:id/traffic", clientHandler.GetTraffic)
			}

			// Node management
			nodeHandler := NewNodeHandler(r.services.Node, r.services.Audit, r.log)
			nodes := protected.Group("/nodes")
			{
				nodes.GET("", nodeHandler.List)
				nodes.GET("/:id", nodeHandler.GetByID)
				nodes.POST("", middleware.RequireRole("admin"), nodeHandler.Create)
				nodes.PUT("/:id", middleware.RequireRole("admin"), nodeHandler.Update)
				nodes.DELETE("/:id", middleware.RequireRole("admin"), nodeHandler.Delete)
			}

			// Traffic logs
			trafficHandler := NewTrafficHandler(r.services.Traffic, r.log)
			protected.GET("/traffic/clients/:id", trafficHandler.GetClientTraffic)
			protected.GET("/traffic/nodes/:id", trafficHandler.GetNodeTraffic)

			// Audit logs
			auditHandler := NewAuditHandler(r.services.Audit, r.log)
			protected.GET("/audit", middleware.RequireRole("admin"), auditHandler.List)

			// 2FA management
			twofa := protected.Group("/2fa")
			{
				twofa.POST("/setup", authHandler.Setup2FA)
				twofa.POST("/enable", authHandler.Enable2FA)
				twofa.POST("/disable", authHandler.Disable2FA)
			}

			// Metrics summary
			protected.GET("/metrics/summary", NewMetricsHandler(r.services, r.log).GetSummary)

			// Inbound management
			inboundHandler := NewInboundHandler(r.services.Inbound, r.services.Audit, r.log)
			inbounds := protected.Group("/inbounds")
			{
				inbounds.GET("", inboundHandler.ListInbounds)
				inbounds.GET("/:id", inboundHandler.GetInbound)
				inbounds.POST("", middleware.RequireRole("admin"), inboundHandler.CreateInbound)
				inbounds.PUT("/:id", middleware.RequireRole("admin"), inboundHandler.UpdateInbound)
				inbounds.DELETE("/:id", middleware.RequireRole("admin"), inboundHandler.DeleteInbound)
				inbounds.POST("/:id/clients", middleware.RequireRole("admin", "operator"), inboundHandler.AddClient)
				inbounds.DELETE("/:id/clients/:email", middleware.RequireRole("admin", "operator"), inboundHandler.RemoveClient)
				inbounds.GET("/:id/clients/:email/link", inboundHandler.GetShareLink)
				inbounds.POST("/:id/toggle", middleware.RequireRole("admin"), inboundHandler.ToggleInbound)
				inbounds.GET("/:id/stats", inboundHandler.GetInboundStats)
			}

			// Template management
			templateHandler := NewTemplateHandler(r.services.Template, r.services.Audit, r.log)
			templates := protected.Group("/templates")
			{
				templates.GET("", templateHandler.ListTemplates)
				templates.GET("/defaults", templateHandler.GetDefaultTemplates)
				templates.GET("/:id", templateHandler.GetTemplate)
				templates.POST("", middleware.RequireRole("admin"), templateHandler.CreateTemplate)
				templates.PUT("/:id", middleware.RequireRole("admin"), templateHandler.UpdateTemplate)
				templates.DELETE("/:id", middleware.RequireRole("admin"), templateHandler.DeleteTemplate)
				templates.POST("/:id/apply", middleware.RequireRole("admin", "operator"), templateHandler.ApplyTemplate)
			}

			// Backup management
			backupHandler := NewBackupHandler(r.services.Backup, r.services.Audit, r.log)
			backup := protected.Group("/backup")
			{
				backup.POST("/create", middleware.RequireRole("admin"), backupHandler.CreateBackup)
				backup.GET("/list", middleware.RequireRole("admin"), backupHandler.ListBackups)
				backup.POST("/restore", middleware.RequireRole("admin"), backupHandler.RestoreBackup)
				backup.DELETE("/:filename", middleware.RequireRole("admin"), backupHandler.DeleteBackup)
			}
		}

		// Public subscription endpoint (no auth required)
		subscriptionHandler := NewSubscriptionHandler(r.services.Subscription, r.log)
		v1.GET("/sub/:token", subscriptionHandler.GetSubscription)
		
		// Protected subscription endpoints
		protected.GET("/subscription/token", subscriptionHandler.GenerateToken)
		protected.GET("/subscription/clients", subscriptionHandler.GetUserClients)
	}
}

