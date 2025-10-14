package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/auth"
	"github.com/v-ui/backend/internal/config"
	"github.com/v-ui/backend/internal/db"
	"github.com/v-ui/backend/internal/models"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user account is disabled")
	Err2FARequired        = errors.New("2FA verification required")
	ErrInvalid2FACode     = errors.New("invalid 2FA code")
)

type AuthService struct {
	db         *db.Database
	jwtManager *auth.JWTManager
	cfg        *config.Config
	log        *zap.Logger
}

func NewAuthService(database *db.Database, jwtManager *auth.JWTManager, cfg *config.Config, log *zap.Logger) *AuthService {
	return &AuthService{
		db:         database,
		jwtManager: jwtManager,
		cfg:        cfg,
		log:        log,
	}
}

func (s *AuthService) Login(username, password string) (*auth.TokenPair, bool, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, false, ErrInvalidCredentials
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, false, ErrInvalidCredentials
	}

	if !user.Enabled {
		return nil, false, ErrUserDisabled
	}

	// If 2FA is enabled, return that it's required
	if user.TwoFactorEnabled {
		return nil, true, Err2FARequired
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(&user)
	if err != nil {
		return nil, false, err
	}

	// Store refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(s.cfg.Auth.RefreshExpiry),
	}
	if err := s.db.Create(refreshToken).Error; err != nil {
		s.log.Error("failed to store refresh token", zap.Error(err))
	}

	return tokens, false, nil
}

func (s *AuthService) Verify2FA(username, code string) (*auth.TokenPair, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.TwoFactorEnabled {
		return nil, errors.New("2FA is not enabled for this user")
	}

	if !auth.ValidateTOTP(user.TwoFactorSecret, code) {
		return nil, ErrInvalid2FACode
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(&user)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(s.cfg.Auth.RefreshExpiry),
	}
	if err := s.db.Create(refreshToken).Error; err != nil {
		s.log.Error("failed to store refresh token", zap.Error(err))
	}

	return tokens, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*auth.TokenPair, error) {
	userID, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if token exists and not revoked
	var storedToken models.RefreshToken
	if err := s.db.Where("token = ? AND user_id = ? AND revoked_at IS NULL", refreshToken, userID).First(&storedToken).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if !user.Enabled {
		return nil, ErrUserDisabled
	}

	// Revoke old token
	now := time.Now()
	storedToken.RevokedAt = &now
	s.db.Save(&storedToken)

	// Generate new tokens
	tokens, err := s.jwtManager.GenerateTokenPair(&user)
	if err != nil {
		return nil, err
	}

	// Store new refresh token
	newRefreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		ExpiresAt: time.Now().Add(s.cfg.Auth.RefreshExpiry),
	}
	if err := s.db.Create(newRefreshToken).Error; err != nil {
		s.log.Error("failed to store refresh token", zap.Error(err))
	}

	return tokens, nil
}

func (s *AuthService) Setup2FA(userID uuid.UUID) (string, string, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return "", "", errors.New("user not found")
	}

	secret, url, err := auth.GenerateTOTPSecret(user.Email, "3X-UI Modern")
	if err != nil {
		return "", "", err
	}

	user.TwoFactorSecret = secret
	if err := s.db.Save(&user).Error; err != nil {
		return "", "", err
	}

	return secret, url, nil
}

func (s *AuthService) Enable2FA(userID uuid.UUID, code string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	if !auth.ValidateTOTP(user.TwoFactorSecret, code) {
		return ErrInvalid2FACode
	}

	user.TwoFactorEnabled = true
	return s.db.Save(&user).Error
}

func (s *AuthService) Disable2FA(userID uuid.UUID) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	user.TwoFactorEnabled = false
	user.TwoFactorSecret = ""
	return s.db.Save(&user).Error
}

