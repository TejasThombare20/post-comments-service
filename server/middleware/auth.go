package middleware

import (
	"net/http"

	"github.com/TejasThombare20/post-comments-service/services"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides authentication middleware using JWT
func AuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from header
		token, err := jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			utils.UnauthorizedResponse(c, err.Error())
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		// Check if it's an access token
		if claims.Type != "access" {
			utils.UnauthorizedResponse(c, "Invalid token type")
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// OptionalAuthMiddleware provides optional authentication
func OptionalAuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Extract token from header
			token, err := jwtService.ExtractTokenFromHeader(authHeader)
			if err == nil {
				// Validate token
				claims, err := jwtService.ValidateToken(token)
				if err == nil && claims.Type == "access" {
					// Set user information in context
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("user_email", claims.Email)
				}
			}
		}
		c.Next()
	}
}

// CORS middleware for handling cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
