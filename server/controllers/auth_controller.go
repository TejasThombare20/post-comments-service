package controllers

import (
	"net/http"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/services"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/TejasThombare20/post-comments-service/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authService services.AuthService
	validator   *validator.Validator
}

// NewAuthController creates a new authentication controller instance
func NewAuthController(authService services.AuthService, validator *validator.Validator) *AuthController {
	return &AuthController{
		authService: authService,
		validator:   validator,
	}
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Validate request
	if validationErrors := ac.validator.ValidateStruct(&req); validationErrors != nil {
		utils.ValidationErrorResponse(c, validationErrors.Error())
		return
	}

	// Register user
	authResponse, err := ac.authService.Register(&req)
	if err != nil {
		if err == utils.ErrUserExists {
			utils.ConflictResponse(c, "User already exists")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to register user")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, authResponse)
}

// Login handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Validate request
	if validationErrors := ac.validator.ValidateStruct(&req); validationErrors != nil {
		utils.ValidationErrorResponse(c, validationErrors.Error())
		return
	}

	// Login user
	authResponse, err := ac.authService.Login(&req)
	if err != nil {
		if err == utils.ErrInvalidCredentials {
			utils.UnauthorizedResponse(c, "Invalid username or password")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to login user")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, authResponse)
}

// RefreshToken handles token refresh
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Validate request
	if validationErrors := ac.validator.ValidateStruct(&req); validationErrors != nil {
		utils.ValidationErrorResponse(c, validationErrors.Error())
		return
	}

	// Refresh token
	authResponse, err := ac.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid refresh token")
		return
	}
	// Use secure response to avoid exposing refresh token
	utils.SuccessResponse(c, http.StatusOK, authResponse.ToSecureResponse())
}

// ChangePassword handles password change requests
func (ac *AuthController) ChangePassword(c *gin.Context) {
	utils.LogInfo("Change password request received", utils.LogFields{})

	// Get user ID from context (set by auth middleware)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogError("User not authenticated for password change", err, utils.LogFields{})
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req models.ChangePasswordRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid request payload for password change", err, utils.LogFields{
			"user_id": userID,
		})
		utils.ValidationErrorResponse(c, "Invalid request format")
		return
	}

	// Validate request
	if validationErrors := ac.validator.ValidateStruct(&req); validationErrors != nil {
		utils.LogError("Validation failed for password change", nil, utils.LogFields{
			"user_id": userID,
			"errors":  validationErrors,
		})
		utils.ValidationErrorResponse(c, validationErrors.Error())
		return
	}

	// Change password
	err = ac.authService.ChangePassword(userID, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			utils.LogError("Invalid current password", err, utils.LogFields{
				"user_id": userID,
			})
			utils.ValidationErrorResponse(c, err.Error())
			return
		}
		if utils.IsNotFoundError(err) {
			utils.LogError("User not found for password change", err, utils.LogFields{
				"user_id": userID,
			})
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.LogError("Failed to change password", err, utils.LogFields{
			"user_id": userID,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("Password changed successfully", utils.LogFields{
		"user_id": userID,
	})

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// Logout handles user logout (client-side token invalidation)
func (ac *AuthController) Logout(c *gin.Context) {
	// in logout call we are removing the access_token from the browser local storage
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Logout successful"})
}

// GetProfile returns the current user's profile
func (ac *AuthController) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID format")
		return
	}

	// This would require access to user service - for now, return basic info from token
	username, _ := c.Get("username")
	email, _ := c.Get("user_email")

	profile := gin.H{
		"id":       userID,
		"username": username,
		"email":    email,
	}

	utils.SuccessResponse(c, http.StatusOK, profile)
}
