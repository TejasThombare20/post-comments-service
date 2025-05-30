package controllers

import (
	"net/http"
	"strconv"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/services"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userService services.UserService
}

// NewUserController creates a new user controller instance
func NewUserController(userService services.UserService) *UserController {
	utils.LogInfo("Initializing user controller", utils.LogFields{
		"component": "user_controller",
	})
	return &UserController{
		userService: userService,
	}
}

// GetUser handles GET /users/:id
func (uc *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	utils.LogInfo("Getting user by ID", utils.LogFields{
		"user_id": idParam,
	})

	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.LogError("Invalid user ID format", err, utils.LogFields{
			"user_id": idParam,
		})
		utils.ValidationErrorResponse(c, "Invalid user ID format")
		return
	}

	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("User not found", err, utils.LogFields{
				"user_id": userID,
			})
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.LogError("Failed to get user", err, utils.LogFields{
			"user_id": userID,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("User retrieved successfully", utils.LogFields{
		"user_id":  userID,
		"username": user.Username,
	})

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// GetUserByUsername handles GET /users/username/:username
func (uc *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	utils.LogInfo("Getting user by username", utils.LogFields{
		"username": username,
	})

	user, err := uc.userService.GetUserByUsername(username)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("User not found by username", err, utils.LogFields{
				"username": username,
			})
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.LogError("Failed to get user by username", err, utils.LogFields{
			"username": username,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("User retrieved by username successfully", utils.LogFields{
		"user_id":  user.ID,
		"username": username,
	})

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// UpdateUser handles PUT /users/:id
func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	utils.LogInfo("Updating user", utils.LogFields{
		"user_id": idParam,
	})

	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.LogError("Invalid user ID format for update", err, utils.LogFields{
			"user_id": idParam,
		})
		utils.ValidationErrorResponse(c, "Invalid user ID format")
		return
	}

	// Get authenticated user ID from context
	authUserID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogError("User not authenticated for update", err, utils.LogFields{
			"target_user_id": userID,
		})
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Check if user is updating their own profile
	if authUserID != userID {
		utils.LogError("User attempting to update another user's profile", nil, utils.LogFields{
			"auth_user_id":   authUserID,
			"target_user_id": userID,
		})
		utils.ForbiddenResponse(c, "You can only update your own profile")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid request payload for user update", err, utils.LogFields{
			"user_id": userID,
		})
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	user, err := uc.userService.UpdateUser(userID, &req)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("User not found for update", err, utils.LogFields{
				"user_id": userID,
			})
			utils.NotFoundResponse(c, "User")
			return
		}
		if utils.IsConflictError(err) {
			utils.LogError("User update conflict", err, utils.LogFields{
				"user_id": userID,
			})
			utils.ConflictResponse(c, "Username or email already exists")
			return
		}
		utils.LogError("Failed to update user", err, utils.LogFields{
			"user_id": userID,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("User updated successfully", utils.LogFields{
		"user_id":    userID,
		"updated_by": authUserID,
	})

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// DeleteUser handles DELETE /users/:id
func (uc *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	utils.LogInfo("Deleting user", utils.LogFields{
		"user_id": idParam,
	})

	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.LogError("Invalid user ID format for deletion", err, utils.LogFields{
			"user_id": idParam,
		})
		utils.ValidationErrorResponse(c, "Invalid user ID format")
		return
	}

	// Get authenticated user ID from context
	authUserID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogError("User not authenticated for deletion", err, utils.LogFields{
			"target_user_id": userID,
		})
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Check if user is deleting their own account
	if authUserID != userID {
		utils.LogError("User attempting to delete another user's account", nil, utils.LogFields{
			"auth_user_id":   authUserID,
			"target_user_id": userID,
		})
		utils.ForbiddenResponse(c, "You can only delete your own account")
		return
	}

	err = uc.userService.DeleteUser(userID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("User not found for deletion", err, utils.LogFields{
				"user_id": userID,
			})
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.LogError("Failed to delete user", err, utils.LogFields{
			"user_id": userID,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("User deleted successfully", utils.LogFields{
		"user_id":    userID,
		"deleted_by": authUserID,
	})

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListUsers handles GET /users
func (uc *UserController) ListUsers(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	utils.LogInfo("Listing users", utils.LogFields{
		"limit":  limitStr,
		"offset": offsetStr,
	})

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		utils.LogError("Invalid limit parameter", err, utils.LogFields{
			"limit": limitStr,
		})
		utils.ValidationErrorResponse(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.LogError("Invalid offset parameter", err, utils.LogFields{
			"offset": offsetStr,
		})
		utils.ValidationErrorResponse(c, "Invalid offset parameter")
		return
	}

	users, err := uc.userService.ListUsers(limit, offset)
	if err != nil {
		utils.LogError("Failed to list users", err, utils.LogFields{
			"limit":  limit,
			"offset": offset,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	// Convert to response format
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	utils.LogInfo("Users listed successfully", utils.LogFields{
		"count":  len(users),
		"limit":  limit,
		"offset": offset,
	})

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"users":  userResponses,
		"limit":  limit,
		"offset": offset,
		"count":  len(userResponses),
	})
}
