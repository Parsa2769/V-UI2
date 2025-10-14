package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/v-ui/backend/internal/api"
	"github.com/v-ui/backend/internal/config"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/logger"
	"github.com/v-ui/backend/internal/metrics"
	"github.com/v-ui/backend/internal/middleware"
	"github.com/v-ui/backend/internal/service"
	"github.com/v-ui/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title V-UI API
// @version 2.0.0
// @description Advanced Xray Control Panel API
// @termsOfService https://github.com/yourusername/v-ui

// @contact.name API Support
// @contact.url https://github.com/yourusername/v-ui/issues
// @contact.email support@example.com

// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Sync()

	log.Info("Starting V-UI",
		zap.String("version", "2.0.0"),
		zap.String("mode", cfg.Server.Mode),
	)

	// Initialize database
	database, err := db.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	log.Info("Database connected successfully")

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	log.Info("Migrations completed successfully")

	// Initialize services
	services := service.NewServices(database, cfg, log)

	// Initialize metrics
	if cfg.Monitoring.EnableMetrics {
		metrics.InitMetrics()
		log.Info("Prometheus metrics enabled", zap.String("path", cfg.Monitoring.MetricsPath))
	}

	// Initialize WebSocket hub
	var wsHub *websocket.Hub
	if cfg.Monitoring.EnableWebSocket {
		wsHub = websocket.NewHub(log)
		go wsHub.Run()
		log.Info("WebSocket hub started", zap.String("path", cfg.Monitoring.WebSocketPath))
	}

	// Set Gin mode
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.CORS(cfg.CORS))
	router.Use(middleware.RequestID())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Unix(),
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		if err := database.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Setup API routes
	apiRouter := api.NewRouter(router, services, cfg, log, wsHub)
	apiRouter.Setup()

	log.Info("API routes configured")

	// Create server
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Info("Server starting", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if wsHub != nil {
		wsHub.Shutdown()
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}

