package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/naputt/cinema-ticket-backend/internal/usecase"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	frontendURL string
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase, frontendURL string) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, frontendURL: frontendURL}
}

// GoogleLogin redirects to Google OAuth consent screen.
// @Summary Redirect to Google OAuth
// @Tags auth
// @Produce json
// @Success 302
// @Router /api/v1/auth/google/login [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.authUsecase.GetAuthURL()
	c.Redirect(http.StatusFound, url)
}

// GoogleCallback handles Google OAuth callback.
// @Summary Handle Google OAuth callback
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization code"})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.HandleCallback(c.Request.Context(), code)
	if err != nil {
		slog.Error("google oauth callback failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	expiry := int(h.authUsecase.GetAccessExpiry().Seconds())
<<<<<<< HEAD
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, expiry, "/", "", h.authUsecase.IsProduction(), true)
=======
	refreshExpiry := int(h.authUsecase.GetRefreshExpiry().Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, expiry, "/", "", h.authUsecase.IsProduction(), true)
	c.SetCookie("refresh_token", refreshToken, refreshExpiry, "/", "", h.authUsecase.IsProduction(), true)
>>>>>>> feature/refresh-cookie

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// RefreshToken refreshes the access token.
// @Summary Refresh access token
// @Tags auth
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshTokenStr, _ := c.Cookie("refresh_token")
	if refreshTokenStr == "" {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh_token"})
			return
		}
		refreshTokenStr = req.RefreshToken
	}

	accessToken, err := h.authUsecase.RefreshAccessToken(c.Request.Context(), refreshTokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	expiry := int(h.authUsecase.GetAccessExpiry().Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", accessToken, expiry, "/", "", h.authUsecase.IsProduction(), true)

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// Logout clears the access_token cookie.
// @Summary Logout
// @Tags auth
// @Success 200 {object} map[string]string
// @Router /api/v1/auth/logout [get]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", "", -1, "/", "", h.authUsecase.IsProduction(), true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// GetMe returns the current user's profile.
// @Summary Get current user
// @Tags auth
// @Produce json
// @Success 200 {object} domain.User
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.authUsecase.GetUserByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
