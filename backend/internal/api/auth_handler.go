package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/v-ui/backend/internal/auth"
	"github.com/v-ui/backend/internal/service"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService  *service.AuthService
	auditService *service.AuditService
	log          *zap.Logger
}

func NewAuthHandler(authService *service.AuthService, auditService *service.AuditService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		auditService: auditService,
		log:          log,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Verify2FARequest struct {
	Username string `json:"username" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Login godoc
// @Summary User login
// @Description Authenticate user and get JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "tokens"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, needs2FA, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == service.Err2FARequired {
			c.JSON(http.StatusOK, gin.H{"requires_2fa": true})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if needs2FA {
		c.JSON(http.StatusOK, gin.H{"requires_2fa": true})
		return
	}

	h.auditService.Log(uuid.Nil, "login", "user:"+req.Username, "User logged in", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, tokens)
}

// Verify2FA godoc
// @Summary Verify 2FA code
// @Description Verify TOTP code and complete login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body Verify2FARequest true "2FA code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/2fa/verify [post]
func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.Verify2FA(req.Username, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(uuid.Nil, "2fa_verify", "user:"+req.Username, "2FA verified", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, tokens)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Setup2FA godoc
// @Summary Setup 2FA
// @Description Generate TOTP secret and QR code for 2FA setup
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /2fa/setup [post]
func (h *AuthHandler) Setup2FA(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	secret, url, err := h.authService.Setup2FA(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "2fa_setup", "self", "2FA setup initiated", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"qr_url": url,
	})
}

// Enable2FA godoc
// @Summary Enable 2FA
// @Description Enable 2FA by verifying a TOTP code
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "TOTP code"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /2fa/enable [post]
func (h *AuthHandler) Enable2FA(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Enable2FA(claims.UserID, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "2fa_enable", "self", "2FA enabled", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FA godoc
// @Summary Disable 2FA
// @Description Disable two-factor authentication
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /2fa/disable [post]
func (h *AuthHandler) Disable2FA(c *gin.Context) {
	claims := c.MustGet("user").(*auth.Claims)

	if err := h.authService.Disable2FA(claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.auditService.Log(claims.UserID, "2fa_disable", "self", "2FA disabled", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

