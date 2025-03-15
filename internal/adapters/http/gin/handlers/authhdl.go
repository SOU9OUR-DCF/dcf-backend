package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/config"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService ports.AuthService
	config      *config.Config
}

func NewAuthHandler(authService ports.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}



func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookie(c, token, time.Now().Add(time.Hour*24*7))
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	res, newToken, err := h.authService.RefreshToken(c.Request.Context(), token)
	if err != nil {
		c.SetCookie("auth_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setAuthCookie(c, newToken, res.ExpiresAt)

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "already logged out"})
		return
	}

	err = h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	secure := h.config.Server.Environment == "prod"
	c.SetCookie("auth_token", "", -1, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}



func (h *AuthHandler) RegisterRestaurant(c *gin.Context) {
	var req domain.RestaurantRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, token, err := h.authService.RegisterRestaurant(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.setAuthCookie(c, token, res.ExpiresAt)
	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) RegisterVolunteer(c *gin.Context) {
	var req domain.VolunteerRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, token, err := h.authService.RegisterVolunteer(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.setAuthCookie(c, token, res.ExpiresAt)
	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) setAuthCookie(c *gin.Context, token domain.Token, expiresAt time.Time) {
	fmt.Println("h.config.Server.Environment", h.config.Server.Environment)
	secure := h.config.Server.Environment == "prod"
	maxAge := int(time.Until(expiresAt).Seconds())
	c.SetCookie(
		"auth_token",
		token.String(),
		maxAge,
		"/",
		"",
		secure,
		true,
	)
}
