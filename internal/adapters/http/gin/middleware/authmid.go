package middleware

import (
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService ports.AuthService
}

func NewAuthMiddleware(authService ports.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(401, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		user, profile, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(401, gin.H{"error": "invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("profile", profile)
		c.Next()
	}
}
