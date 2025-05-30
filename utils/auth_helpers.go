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

// GetUsernameFromContext extracts the username from the Gin context
func GetUsernameFromContext(c *gin.Context) (string, error) {
	usernameInterface, exists := c.Get("username")
	if !exists {
		return "", errors.New("username not found in context")
	}

	username, ok := usernameInterface.(string)
	if !ok {
		return "", errors.New("invalid username format")
	}

	return username, nil
}

// GetUserEmailFromContext extracts the user email from the Gin context
func GetUserEmailFromContext(c *gin.Context) (*string, error) {
	emailInterface, exists := c.Get("user_email")
	if !exists {
		return nil, nil // Email is optional
	}

	if emailInterface == nil {
		return nil, nil
	}

	email, ok := emailInterface.(*string)
	if !ok {
		return nil, errors.New("invalid email format")
	}

	return email, nil
}

// IsUserAuthenticated checks if a user is authenticated in the current context
func IsUserAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}
