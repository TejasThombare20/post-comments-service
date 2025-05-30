package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	StatusCode   int         `json:"status_code"`
	ErrorMessage *string     `json:"error_message"`
	Data         interface{} `json:"data"`
}

// SuccessResponse sends a successful response with data
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		StatusCode:   statusCode,
		ErrorMessage: nil,
		Data:         data,
	})
}

// ErrorResponse sends an error response with error message
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		StatusCode:   statusCode,
		ErrorMessage: &message,
		Data:         nil,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, resource+" not found")
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, message)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	ErrorResponse(c, http.StatusForbidden, message)
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message)
}
