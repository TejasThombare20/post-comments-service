package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        *string    `json:"email" db:"email"`
	PasswordHash *string    `json:"-" db:"password_hash"`
	DisplayName  *string    `json:"display_name" db:"display_name"`
	AvatarURL    *string    `json:"avatar_url" db:"avatar_url"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"-" db:"deleted_at"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Username    string  `json:"username" validate:"required,min=3,max=50"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Password    string  `json:"password" validate:"required,min=6"`
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
	AvatarURL   *string `json:"avatar_url" validate:"omitempty,url"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email       *string `json:"email" validate:"omitempty,email"`
	DisplayName *string `json:"display_name" validate:"omitempty,max=100"`
	AvatarURL   *string `json:"avatar_url" validate:"omitempty,url"`
}

// GetUserRequest represents the request payload for getting a user by ID
type GetUserRequest struct {
	ID string `json:"id" validate:"required,uuid" uri:"id"`
}

// GetUserByUsernameRequest represents the request payload for getting a user by username
type GetUserByUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" uri:"username"`
}

// DeleteUserRequest represents the request payload for deleting a user
type DeleteUserRequest struct {
	ID string `json:"id" validate:"required,uuid" uri:"id"`
}

// ListUsersRequest represents the request payload for listing users
type ListUsersRequest struct {
	Limit  int `json:"limit" validate:"omitempty,gte=1,lte=100" form:"limit"`
	Offset int `json:"offset" validate:"omitempty,gte=0" form:"offset"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       *string   `json:"email"`
	DisplayName *string   `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
