package services

import (
	"errors"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService interface defines authentication business logic methods
type AuthService interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(req *models.RefreshTokenRequest) (*models.AuthResponse, error)
	ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error
	ValidateToken(tokenString string) (*models.JWTClaims, error)
}

// authService implements AuthService interface
type authService struct {
	userService UserService
	jwtService  *JWTService
}

// NewAuthService creates a new authentication service instance
func NewAuthService(userService UserService, jwtService *JWTService) AuthService {
	utils.LogInfo("Initializing authentication service", utils.LogFields{
		"component": "auth_service",
	})
	return &authService{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Register creates a new user account and returns authentication tokens
func (s *authService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	utils.LogInfo("Starting user registration", utils.LogFields{
		"username": req.Username,
		"email":    req.Email,
	})

	// Convert RegisterRequest to CreateUserRequest
	createUserReq := &models.CreateUserRequest{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
	}

	// Create the user
	user, err := s.userService.CreateUser(createUserReq)
	if err != nil {
		utils.LogError("User registration failed", err, utils.LogFields{
			"username": req.Username,
			"email":    req.Email,
		})
		return nil, err
	}

	// Generate JWT tokens
	authResponse, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		utils.LogError("Token generation failed during registration", err, utils.LogFields{
			"user_id":  user.ID,
			"username": req.Username,
		})
		return nil, utils.WrapError(err, "failed to generate tokens")
	}

	utils.LogInfo("User registration completed successfully", utils.LogFields{
		"user_id":  user.ID,
		"username": req.Username,
		"email":    req.Email,
	})

	return authResponse, nil
}

// Login authenticates a user and returns authentication tokens
func (s *authService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	utils.LogInfo("Starting user login", utils.LogFields{
		"username": req.Username,
	})

	// Get user by username
	user, err := s.userService.GetUserByUsername(req.Username)
	if err != nil {
		utils.LogError("User not found during login", err, utils.LogFields{
			"username": req.Username,
		})
		return nil, utils.ErrInvalidCredentials
	}

	// Check if password hash exists
	if user.PasswordHash == nil {
		utils.LogError("User has no password hash", nil, utils.LogFields{
			"user_id":  user.ID,
			"username": req.Username,
		})
		return nil, utils.ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password))
	if err != nil {
		utils.LogError("Password verification failed during login", err, utils.LogFields{
			"user_id":  user.ID,
			"username": req.Username,
		})
		return nil, utils.ErrInvalidCredentials
	}

	// Generate JWT tokens
	authResponse, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		utils.LogError("Token generation failed during login", err, utils.LogFields{
			"user_id":  user.ID,
			"username": req.Username,
		})
		return nil, utils.WrapError(err, "failed to generate tokens")
	}

	utils.LogInfo("User login completed successfully", utils.LogFields{
		"user_id":  user.ID,
		"username": req.Username,
	})

	return authResponse, nil
}

// RefreshToken generates new tokens using a valid refresh token
func (s *authService) RefreshToken(req *models.RefreshTokenRequest) (*models.AuthResponse, error) {
	utils.LogInfo("Starting token refresh", utils.LogFields{
		"refresh_token_length": len(req.RefreshToken),
	})

	authResponse, err := s.jwtService.RefreshToken(req.RefreshToken, s.userService)
	if err != nil {
		utils.LogError("Token refresh failed", err, nil)
		return nil, err
	}

	utils.LogInfo("Token refresh completed successfully", utils.LogFields{
		"user_id": authResponse.User.ID,
	})

	return authResponse, nil
}

// ChangePassword changes a user's password
func (s *authService) ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error {
	utils.LogInfo("Starting password change process", utils.LogFields{
		"user_id": userID,
	})

	// Get the user
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		utils.LogError("User not found for password change", err, utils.LogFields{
			"user_id": userID,
		})
		return err
	}

	// Check if password hash exists
	if user.PasswordHash == nil {
		utils.LogError("User has no password set", nil, utils.LogFields{
			"user_id": userID,
		})
		return errors.New("user has no password set")
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		utils.LogError("Current password verification failed", err, utils.LogFields{
			"user_id": userID,
		})
		return utils.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("Failed to hash new password", err, utils.LogFields{
			"user_id": userID,
		})
		return utils.WrapError(err, "failed to hash new password")
	}

	// Update user password
	err = s.userService.UpdatePassword(userID, string(hashedPassword))
	if err != nil {
		utils.LogError("Failed to update user password", err, utils.LogFields{
			"user_id": userID,
		})
		return utils.WrapError(err, "failed to update password")
	}

	utils.LogInfo("Password changed successfully", utils.LogFields{
		"user_id": userID,
	})

	return nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *authService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	return s.jwtService.ValidateToken(tokenString)
}
