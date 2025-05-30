package utils

import (
	"errors"
	"fmt"
)

// Custom error types
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrPostNotFound          = errors.New("post not found")
	ErrCommentNotFound       = errors.New("comment not found")
	ErrUserExists            = errors.New("user already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUnauthorized          = errors.New("unauthorized access")
	ErrForbidden             = errors.New("forbidden access")
	ErrInvalidInput          = errors.New("invalid input")
	ErrDatabaseError         = errors.New("database error")
	ErrInternalServer        = errors.New("internal server error")
)

// ErrorMessages contains predefined error messages for different scenarios
type ErrorMessages struct {
	UserNotFound        string
	PostNotFound        string
	CommentNotFound     string
	UserAlreadyExists   string
	InvalidCredentials  string
	Unauthorized        string
	Forbidden           string
	InvalidInput        string
	DatabaseError       string
	InternalServerError string
	ValidationFailed    string
}

// GetErrorMessages returns predefined error messages
func GetErrorMessages() ErrorMessages {
	return ErrorMessages{
		UserNotFound:        "User not found",
		PostNotFound:        "Post not found",
		CommentNotFound:     "Comment not found",
		UserAlreadyExists:   "User with this username or email already exists",
		InvalidCredentials:  "Invalid username or password",
		Unauthorized:        "You are not authorized to perform this action",
		Forbidden:           "Access to this resource is forbidden",
		InvalidInput:        "Invalid input provided",
		DatabaseError:       "Database operation failed",
		InternalServerError: "An internal server error occurred",
		ValidationFailed:    "Validation failed",
	}
}

// WrapError wraps an error with additional context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrPostNotFound) ||
		errors.Is(err, ErrCommentNotFound)
}

// IsConflictError checks if the error is a conflict error
func IsConflictError(err error) bool {
	return errors.Is(err, ErrUserExists) ||
		errors.Is(err, ErrUsernameAlreadyExists) ||
		errors.Is(err, ErrEmailAlreadyExists)
}

// IsUnauthorizedError checks if the error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	return errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrInvalidCredentials)
}

// IsForbiddenError checks if the error is a forbidden error
func IsForbiddenError(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}
