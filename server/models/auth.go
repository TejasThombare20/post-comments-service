package models

import (
	"time"

	"github.com/google/uuid"
)

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=50"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Password    string  `json:"password" validate:"required,min=6"`
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
	AvatarURL   *string `json:"avatar_url" validate:"omitempty,url"`
}

// AuthResponse represents the response payload for authentication
// NOTE: This includes refresh_token and should only be used for login/register endpoints
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// SecureAuthResponse represents the secure response payload for authentication
// This excludes refresh_token for better security
type SecureAuthResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token"`
	ExpiresAt   time.Time    `json:"expires_at"`
}

// ToSecureResponse converts AuthResponse to SecureAuthResponse
func (a *AuthResponse) ToSecureResponse() SecureAuthResponse {
	return SecureAuthResponse{
		User:        a.User,
		AccessToken: a.AccessToken,
		ExpiresAt:   a.ExpiresAt,
	}
}

// RefreshTokenRequest represents the request payload for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    *string   `json:"email"`
	Type     string    `json:"type"` // "access" or "refresh"
}

// ChangePasswordRequest represents the request payload for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}
