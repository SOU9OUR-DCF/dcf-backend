package middleware

import (
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/config"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := cfg.CORS.AllowedOrigins
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{"*"}
		}

		origin := c.Request.Header.Get("Origin")
		allowOrigin := "*"

		if len(allowedOrigins) > 0 && allowedOrigins[0] != "*" {
			allowOrigin = ""
			for _, allowed := range allowedOrigins {
				if allowed == origin {
					allowOrigin = origin
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
