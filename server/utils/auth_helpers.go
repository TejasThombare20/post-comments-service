package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromContext extracts the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user not authenticated")
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	return userID, nil
}
