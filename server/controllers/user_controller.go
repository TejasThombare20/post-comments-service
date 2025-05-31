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
	return &UserController{
		userService: userService,
	}
}

// GetUserByID handles GET /users/:id
func (uc *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID format")
		return
	}

	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// GetUserByUsername handles GET /users/username/:username
func (uc *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.ValidationErrorResponse(c, "Username is required")
		return
	}

	user, err := uc.userService.GetUserByUsername(username)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// UpdateUser handles PUT /users/:id
func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID format")
		return
	}

	authenticatedUserID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if userID != authenticatedUserID {
		utils.ForbiddenResponse(c, "You can only update your own profile")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	user, err := uc.userService.UpdateUser(userID, &req)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "User")
			return
		}
		if utils.IsConflictError(err) {
			utils.ConflictResponse(c, "Email already exists")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user.ToResponse())
}

// DeleteUser handles DELETE /users/:id
func (uc *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID format")
		return
	}

	authenticatedUserID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if userID != authenticatedUserID {
		utils.ForbiddenResponse(c, "You can only delete your own account")
		return
	}

	err = uc.userService.DeleteUser(userID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListUsers handles GET /users
func (uc *UserController) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		utils.ValidationErrorResponse(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.ValidationErrorResponse(c, "Invalid offset parameter")
		return
	}

	users, err := uc.userService.ListUsers(limit, offset)
	if err != nil {
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"users":  userResponses,
		"limit":  limit,
		"offset": offset,
		"count":  len(userResponses),
	})
}
