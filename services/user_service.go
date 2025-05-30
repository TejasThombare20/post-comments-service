package services

import (
	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/repository"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService interface defines user business logic methods
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error)
	UpdatePassword(id uuid.UUID, hashedPassword string) error
	DeleteUser(id uuid.UUID) error
	ListUsers(limit, offset int) ([]models.User, error)
}

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repository.UserRepository) UserService {
	utils.LogInfo("Initializing user service", utils.LogFields{
		"component": "user_service",
	})
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with hashed password
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	utils.LogInfo("Creating new user", utils.LogFields{
		"username": req.Username,
		"email": func() string {
			if req.Email != nil {
				return *req.Email
			}
			return ""
		}(),
	})

	// Check if user already exists by username
	if existingUser, err := s.userRepo.GetByUsername(req.Username); err == nil && existingUser != nil {
		utils.LogError("User already exists with username", nil, utils.LogFields{
			"username": req.Username,
		})
		return nil, utils.ErrUserExists
	}

	// Check if user already exists by email (if provided)
	if req.Email != nil && *req.Email != "" {
		if existingUser, err := s.userRepo.GetByEmail(*req.Email); err == nil && existingUser != nil {
			utils.LogError("User already exists with email", nil, utils.LogFields{
				"email": *req.Email,
			})
			return nil, utils.ErrUserExists
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("Failed to hash password", err, utils.LogFields{
			"username": req.Username,
		})
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &[]string{string(hashedPassword)}[0],
		DisplayName:  req.DisplayName,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		utils.LogError("Failed to create user in database", err, utils.LogFields{
			"username": req.Username,
		})
		return nil, err
	}

	utils.LogInfo("User created successfully", utils.LogFields{
		"user_id":  user.ID,
		"username": user.Username,
	})

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	utils.LogInfo("Getting user by ID", utils.LogFields{
		"user_id": id,
	})

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("User not found", err, utils.LogFields{
				"user_id": id,
			})
		} else {
			utils.LogError("Failed to get user by ID", err, utils.LogFields{
				"user_id": id,
			})
		}
		return nil, err
	}

	utils.LogInfo("User retrieved successfully", utils.LogFields{
		"user_id":  user.ID,
		"username": user.Username,
	})

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	utils.LogInfo("Getting user by username", utils.LogFields{
		"username": username,
	})

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("User not found by username", err, utils.LogFields{
				"username": username,
			})
		} else {
			utils.LogError("Failed to get user by username", err, utils.LogFields{
				"username": username,
			})
		}
		return nil, err
	}

	utils.LogInfo("User retrieved by username successfully", utils.LogFields{
		"user_id":  user.ID,
		"username": username,
	})

	return user, nil
}

// UpdateUser updates user information
func (s *userService) UpdateUser(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	utils.LogInfo("Updating user", utils.LogFields{
		"user_id": id,
	})

	// Check if user exists
	if _, err := s.userRepo.GetByID(id); err != nil {
		utils.LogError("Failed to get user for update", err, utils.LogFields{
			"user_id": id,
		})
		return nil, err
	}

	// Check if email is being updated and already exists
	if req.Email != nil && *req.Email != "" {
		if existingUser, err := s.userRepo.GetByEmail(*req.Email); err == nil && existingUser != nil && existingUser.ID != id {
			utils.LogError("Email already taken by another user", nil, utils.LogFields{
				"user_id":          id,
				"new_email":        *req.Email,
				"existing_user_id": existingUser.ID,
			})
			return nil, utils.ErrUserExists
		}
	}

	// Update user
	updatedUser, err := s.userRepo.Update(id, req)
	if err != nil {
		utils.LogError("Failed to update user in database", err, utils.LogFields{
			"user_id": id,
		})
		return nil, err
	}

	utils.LogInfo("User updated successfully", utils.LogFields{
		"user_id":  updatedUser.ID,
		"username": updatedUser.Username,
	})

	return updatedUser, nil
}

// UpdatePassword updates the password for a user
func (s *userService) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	utils.LogInfo("Updating user password", utils.LogFields{
		"user_id": id,
	})

	// Check if user exists
	if _, err := s.userRepo.GetByID(id); err != nil {
		utils.LogError("Failed to get user for password update", err, utils.LogFields{
			"user_id": id,
		})
		return err
	}

	// Update password
	if err := s.userRepo.UpdatePassword(id, hashedPassword); err != nil {
		utils.LogError("Failed to update password in database", err, utils.LogFields{
			"user_id": id,
		})
		return err
	}

	utils.LogInfo("User password updated successfully", utils.LogFields{
		"user_id": id,
	})

	return nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uuid.UUID) error {
	utils.LogInfo("Deleting user", utils.LogFields{
		"user_id": id,
	})

	// Check if user exists
	if _, err := s.userRepo.GetByID(id); err != nil {
		utils.LogError("Failed to get user for deletion", err, utils.LogFields{
			"user_id": id,
		})
		return err
	}

	// Delete user
	if err := s.userRepo.Delete(id); err != nil {
		utils.LogError("Failed to delete user from database", err, utils.LogFields{
			"user_id": id,
		})
		return err
	}

	utils.LogInfo("User deleted successfully", utils.LogFields{
		"user_id": id,
	})

	return nil
}

// ListUsers retrieves a paginated list of users
func (s *userService) ListUsers(limit, offset int) ([]models.User, error) {
	utils.LogInfo("Listing users", utils.LogFields{
		"limit":  limit,
		"offset": offset,
	})

	// Set default and maximum limits
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.List(limit, offset)
	if err != nil {
		utils.LogError("Failed to list users", err, utils.LogFields{
			"limit":  limit,
			"offset": offset,
		})
		return nil, err
	}

	utils.LogInfo("Users listed successfully", utils.LogFields{
		"count":  len(users),
		"limit":  limit,
		"offset": offset,
	})

	return users, nil
}

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}
