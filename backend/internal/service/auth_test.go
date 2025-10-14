package service

import (
	"testing"
	"time"

	"github.com/v-ui/backend/internal/auth"
	"github.com/v-ui/backend/internal/config"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *db.Database {
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	database := &db.Database{DB: gormDB}
	if err := database.Migrate(); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return database
}

func TestAuthService_Login(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-32-characters-long!",
			JWTRefreshSecret: "test-refresh-32-characters-long",
			AccessExpiry:     24 * time.Hour,
			RefreshExpiry:    168 * time.Hour,
			BcryptCost:       4, // Lower cost for testing
		},
	}

	log, _ := zap.NewDevelopment()
	jwtManager := auth.NewJWTManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTRefreshSecret,
		cfg.Auth.AccessExpiry,
		cfg.Auth.RefreshExpiry,
	)

	authService := NewAuthService(database, jwtManager, cfg, log)
	userService := NewUserService(database, cfg, log)

	// Create test user
	_, err := userService.Create("testuser", "test@example.com", "password123", models.RoleAdmin)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	tests := []struct {
		name        string
		username    string
		password    string
		wantErr     bool
		wantTokens  bool
		wants2FA    bool
	}{
		{
			name:       "valid credentials",
			username:   "testuser",
			password:   "password123",
			wantErr:    false,
			wantTokens: true,
			wants2FA:   false,
		},
		{
			name:       "invalid password",
			username:   "testuser",
			password:   "wrongpassword",
			wantErr:    true,
			wantTokens: false,
			wants2FA:   false,
		},
		{
			name:       "invalid username",
			username:   "nonexistent",
			password:   "password123",
			wantErr:    true,
			wantTokens: false,
			wants2FA:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, needs2FA, err := authService.Login(tt.username, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (tokens != nil) != tt.wantTokens {
				t.Errorf("Login() tokens = %v, wantTokens %v", tokens != nil, tt.wantTokens)
			}

			if needs2FA != tt.wants2FA {
				t.Errorf("Login() needs2FA = %v, wants2FA %v", needs2FA, tt.wants2FA)
			}

			if tokens != nil {
				if tokens.AccessToken == "" {
					t.Error("Login() AccessToken is empty")
				}
				if tokens.RefreshToken == "" {
					t.Error("Login() RefreshToken is empty")
				}
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	database := setupTestDB(t)
	defer database.Close()

	cfg := &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:        "test-secret-32-characters-long!",
			JWTRefreshSecret: "test-refresh-32-characters-long",
			AccessExpiry:     24 * time.Hour,
			RefreshExpiry:    168 * time.Hour,
			BcryptCost:       4,
		},
	}

	log, _ := zap.NewDevelopment()
	jwtManager := auth.NewJWTManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTRefreshSecret,
		cfg.Auth.AccessExpiry,
		cfg.Auth.RefreshExpiry,
	)

	authService := NewAuthService(database, jwtManager, cfg, log)
	userService := NewUserService(database, cfg, log)

	// Create test user
	user, err := userService.Create("testuser", "test@example.com", "password123", models.RoleAdmin)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Login to get tokens
	tokens, _, err := authService.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}

	// Test refresh token
	newTokens, err := authService.RefreshToken(tokens.RefreshToken)
	if err != nil {
		t.Errorf("RefreshToken() error = %v", err)
		return
	}

	if newTokens.AccessToken == "" {
		t.Error("RefreshToken() AccessToken is empty")
	}

	if newTokens.AccessToken == tokens.AccessToken {
		t.Error("RefreshToken() returned same access token")
	}

	// Verify new access token is valid
	claims, err := jwtManager.ValidateAccessToken(newTokens.AccessToken)
	if err != nil {
		t.Errorf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("ValidateAccessToken() UserID = %v, want %v", claims.UserID, user.ID)
	}
}

