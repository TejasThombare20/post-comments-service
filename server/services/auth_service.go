package services

import (
	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/repository"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/TejasThombare20/post-comments-service/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService interface defines authentication business logic methods
type AuthService interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(refreshToken string) (*models.AuthResponse, error)
	ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error
}

// authService implements AuthService interface
type authService struct {
	userRepo    repository.UserRepository
	jwtService  *JWTService
	userService UserService
	validator   *validator.Validator
}

// NewAuthService creates a new authentication service instance
func NewAuthService(userRepo repository.UserRepository, jwtService *JWTService, userService UserService, validator *validator.Validator) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtService:  jwtService,
		userService: userService,
		validator:   validator,
	}
}

// Register creates a new user account and returns authentication tokens
func (s *authService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	createUserReq := &models.CreateUserRequest{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	}

	user, err := s.userService.CreateUser(createUserReq)
	if err != nil {
		return nil, err
	}

	authResponse, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, utils.WrapError(err, "failed to generate tokens")
	}

	return authResponse, nil
}

// Login authenticates a user and returns authentication tokens
func (s *authService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	if user.PasswordHash == nil {
		return nil, utils.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	authResponse, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, utils.WrapError(err, "failed to generate tokens")
	}

	return authResponse, nil
}

// RefreshToken generates new tokens using a valid refresh token
func (s *authService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	authResponse, err := s.jwtService.RefreshToken(refreshToken, s.userService)
	if err != nil {
		return nil, err
	}

	return authResponse, nil
}

// ChangePassword changes a user's password
func (s *authService) ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error {
	if err := s.validator.ValidateStruct(req); err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if user.PasswordHash == nil {
		return utils.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return utils.ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.WrapError(err, "failed to hash new password")
	}

	if err := s.userService.UpdatePassword(userID, string(hashedPassword)); err != nil {
		return utils.WrapError(err, "failed to update password")
	}

	return nil
}
