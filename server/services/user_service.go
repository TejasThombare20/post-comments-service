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
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with hashed password
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	if existingUser, err := s.userRepo.GetByUsername(req.Username); err == nil && existingUser != nil {
		return nil, utils.ErrUserExists
	}

	if req.Email != nil && *req.Email != "" {
		if existingUser, err := s.userRepo.GetByEmail(*req.Email); err == nil && existingUser != nil {
			return nil, utils.ErrUserExists
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &[]string{string(hashedPassword)}[0],
		DisplayName:  req.DisplayName,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates user information
func (s *userService) UpdateUser(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	if _, err := s.userRepo.GetByID(id); err != nil {
		return nil, err
	}

	if req.Email != nil && *req.Email != "" {
		if existingUser, err := s.userRepo.GetByEmail(*req.Email); err == nil && existingUser != nil && existingUser.ID != id {
			return nil, utils.ErrUserExists
		}
	}

	updatedUser, err := s.userRepo.Update(id, req)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// UpdatePassword updates the password for a user
func (s *userService) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	if _, err := s.userRepo.GetByID(id); err != nil {
		return err
	}

	err := s.userRepo.UpdatePassword(id, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uuid.UUID) error {
	if _, err := s.userRepo.GetByID(id); err != nil {
		return err
	}

	err := s.userRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// ListUsers retrieves a paginated list of users
func (s *userService) ListUsers(limit, offset int) ([]models.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if limit > 100 {
		limit = 100
	}

	users, err := s.userRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}
