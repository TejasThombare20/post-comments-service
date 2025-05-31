package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/google/uuid"
)

// UserRepository interface defines user data access methods
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id uuid.UUID, updates *models.UpdateUserRequest) (*models.User, error)
	UpdatePassword(id uuid.UUID, hashedPassword string) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]models.User, error)
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user in the database
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, display_name, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.AvatarURL,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return utils.WrapError(err, "failed to create user")
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, created_at, updated_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrUserNotFound
		}
		return nil, utils.WrapError(err, "failed to get user by ID")
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, created_at, updated_at
		FROM users 
		WHERE username = $1 AND deleted_at IS NULL`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrUserNotFound
		}
		return nil, utils.WrapError(err, "failed to get user by username")
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, created_at, updated_at
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrUserNotFound
		}
		return nil, utils.WrapError(err, "failed to get user by email")
	}

	return &user, nil
}

// Update updates a user's information
func (r *userRepository) Update(id uuid.UUID, updates *models.UpdateUserRequest) (*models.User, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, updates.Email)
		argIndex++
	}

	if updates.DisplayName != nil {
		setParts = append(setParts, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, updates.DisplayName)
		argIndex++
	}

	if updates.AvatarURL != nil {
		setParts = append(setParts, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, updates.AvatarURL)
		argIndex++
	}

	if len(setParts) == 0 {
		// No updates to perform, just return the current user
		return r.GetByID(id)
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE users 
		SET %s
		WHERE id = $%d AND deleted_at IS NULL`,
		strings.Join(setParts, ", "),
		argIndex,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, utils.WrapError(err, "failed to update user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return nil, utils.ErrUserNotFound
	}

	// Return updated user
	return r.GetByID(id)
}

// UpdatePassword updates a user's password
func (r *userRepository) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	query := `
		UPDATE users 
		SET password_hash = $1 
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, hashedPassword, id)
	if err != nil {
		return utils.WrapError(err, "failed to update user password")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE users 
		SET deleted_at = $1 
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return utils.WrapError(err, "failed to delete user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// List retrieves a paginated list of users
func (r *userRepository) List(limit, offset int) ([]models.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, created_at, updated_at
		FROM users 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list users")
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.DisplayName,
			&user.AvatarURL,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, utils.WrapError(err, "failed to scan user row")
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.WrapError(err, "error iterating user rows")
	}

	return users, nil
}
