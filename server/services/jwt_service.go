package services

import (
	"errors"
	"os"
	"time"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewJWTService creates a new JWT service instance
func NewJWTService() *JWTService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-this-in-production" // Default for development
		utils.LogWarn("Using default JWT secret key - change this in production", nil)
	}

	utils.LogInfo("Initializing JWT service", utils.LogFields{
		"component":         "jwt_service",
		"access_token_ttl":  "15m",
		"refresh_token_ttl": "7d",
	})

	return &JWTService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  15 * time.Minute,   // Access token expires in 15 minutes
		refreshTokenTTL: 7 * 24 * time.Hour, // Refresh token expires in 7 days
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTService) GenerateTokenPair(user *models.User) (*models.AuthResponse, error) {
	utils.LogInfo("Generating token pair", utils.LogFields{
		"user_id":  user.ID,
		"username": user.Username,
	})

	// Generate access token
	accessToken, accessExpiresAt, err := j.generateToken(user, "access", j.accessTokenTTL)
	if err != nil {
		utils.LogError("Failed to generate access token", err, utils.LogFields{
			"user_id": user.ID,
		})
		return nil, err
	}

	// Generate refresh token
	refreshToken, _, err := j.generateToken(user, "refresh", j.refreshTokenTTL)
	if err != nil {
		utils.LogError("Failed to generate refresh token", err, utils.LogFields{
			"user_id": user.ID,
		})
		return nil, err
	}

	utils.LogInfo("Token pair generated successfully", utils.LogFields{
		"user_id":           user.ID,
		"access_expires_at": accessExpiresAt,
	})

	return &models.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt,
	}, nil
}

// generateToken creates a JWT token with the specified type and TTL
func (j *JWTService) generateToken(user *models.User, tokenType string, ttl time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(ttl)

	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"type":     tokenType,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	utils.LogInfo("Validating JWT token", utils.LogFields{
		"token_length": len(tokenString),
	})

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		utils.LogError("JWT token parsing failed", err, nil)
		return nil, err
	}

	if !token.Valid {
		utils.LogError("JWT token is invalid", nil, nil)
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.LogError("JWT token claims are invalid", nil, nil)
		return nil, errors.New("invalid token claims")
	}

	// Parse user ID
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		utils.LogError("Invalid user_id in JWT token", nil, nil)
		return nil, errors.New("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.LogError("Invalid user_id format in JWT token", err, utils.LogFields{
			"user_id_str": userIDStr,
		})
		return nil, errors.New("invalid user_id format")
	}

	// Parse username
	username, ok := claims["username"].(string)
	if !ok {
		utils.LogError("Invalid username in JWT token", nil, nil)
		return nil, errors.New("invalid username in token")
	}

	// Parse email (optional)
	var email *string
	if emailVal, exists := claims["email"]; exists && emailVal != nil {
		if emailStr, ok := emailVal.(string); ok {
			email = &emailStr
		}
	}

	// Parse token type
	tokenType, ok := claims["type"].(string)
	if !ok {
		utils.LogError("Invalid token type in JWT token", nil, nil)
		return nil, errors.New("invalid token type")
	}

	utils.LogInfo("JWT token validated successfully", utils.LogFields{
		"user_id":    userID,
		"username":   username,
		"token_type": tokenType,
	})

	return &models.JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Type:     tokenType,
	}, nil
}

// RefreshToken generates a new access token using a valid refresh token
func (j *JWTService) RefreshToken(refreshTokenString string, userService UserService) (*models.AuthResponse, error) {
	utils.LogInfo("Starting token refresh process", utils.LogFields{
		"refresh_token_length": len(refreshTokenString),
	})

	// Validate the refresh token
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		utils.LogError("Refresh token validation failed", err, nil)
		return nil, err
	}

	// Check if it's a refresh token
	if claims.Type != "refresh" {
		utils.LogError("Invalid token type for refresh", nil, utils.LogFields{
			"token_type": claims.Type,
		})
		return nil, errors.New("invalid token type for refresh")
	}

	// Get the user from database to ensure they still exist
	user, err := userService.GetUserByID(claims.UserID)
	if err != nil {
		utils.LogError("User not found during token refresh", err, utils.LogFields{
			"user_id": claims.UserID,
		})
		return nil, errors.New("user not found")
	}

	// Generate new token pair
	authResponse, err := j.GenerateTokenPair(user)
	if err != nil {
		utils.LogError("Failed to generate new token pair during refresh", err, utils.LogFields{
			"user_id": claims.UserID,
		})
		return nil, err
	}

	utils.LogInfo("Token refresh completed successfully", utils.LogFields{
		"user_id": claims.UserID,
	})

	return authResponse, nil
}

// ExtractTokenFromHeader extracts the token from Authorization header
func (j *JWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization format")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("token required")
	}

	return token, nil
}
