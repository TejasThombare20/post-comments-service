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

// PostRepository interface defines post data access methods
type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uuid.UUID) (*models.Post, error)
	GetByIDWithAuthor(id uuid.UUID) (*models.Post, error)
	GetByIDWithComments(id uuid.UUID) (*models.Post, error)
	Update(id uuid.UUID, updates *models.UpdatePostRequest) (*models.Post, error)
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]models.Post, error)
	ListByUser(userID uuid.UUID, limit, offset int) ([]models.Post, error)
}

// postRepository implements PostRepository interface
type postRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new post repository instance
func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

// Create creates a new post in the database
func (r *postRepository) Create(post *models.Post) error {
	query := `
		INSERT INTO posts (id, title, content, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query,
		post.ID,
		post.Title,
		post.Content,
		post.CreatedBy,
		post.CreatedAt,
		post.UpdatedAt,
	)

	if err != nil {
		return utils.WrapError(err, "failed to create post")
	}

	return nil
}

// GetByID retrieves a post by ID
func (r *postRepository) GetByID(id uuid.UUID) (*models.Post, error) {
	query := `
		SELECT id, title, content, created_by, created_at, updated_at
		FROM posts 
		WHERE id = $1 AND deleted_at IS NULL`

	var post models.Post
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrPostNotFound
		}
		return nil, utils.WrapError(err, "failed to get post by ID")
	}

	return &post, nil
}

// GetByIDWithAuthor retrieves a post by ID with author information
func (r *postRepository) GetByIDWithAuthor(id uuid.UUID) (*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_by, p.created_at, p.updated_at,
		       u.id, u.username, u.email, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.created_by = u.id
		WHERE p.id = $1 AND p.deleted_at IS NULL AND u.deleted_at IS NULL`

	var post models.Post
	var author models.User

	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedBy,
		&post.CreatedAt,
		&post.UpdatedAt,
		&author.ID,
		&author.Username,
		&author.Email,
		&author.DisplayName,
		&author.AvatarURL,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrPostNotFound
		}
		return nil, utils.WrapError(err, "failed to get post with author")
	}

	post.Author = &author
	return &post, nil
}

// GetByIDWithComments retrieves a post by ID with comments and authors
func (r *postRepository) GetByIDWithComments(id uuid.UUID) (*models.Post, error) {
	post, err := r.GetByIDWithAuthor(id)
	if err != nil {
		return nil, err
	}

	commentsQuery := `
		SELECT c.id, c.content, c.post_id, c.parent_id, c.thread_id, c.created_by, c.created_at,
		       u.id, u.username, u.email, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM comments c
		LEFT JOIN users u ON c.created_by = u.id AND u.deleted_at IS NULL
		WHERE c.post_id = $1 AND c.deleted_at IS NULL AND c.parent_id IS NULL
		ORDER BY c.created_at DESC`

	rows, err := r.db.Query(commentsQuery, id)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get comments for post")
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var author models.User
		var authorID sql.NullString
		var authorUsername sql.NullString
		var authorEmail sql.NullString
		var authorDisplayName sql.NullString
		var authorAvatarURL sql.NullString
		var authorCreatedAt sql.NullTime
		var authorUpdatedAt sql.NullTime

		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.ParentID,
			&comment.ThreadID,
			&comment.CreatedBy,
			&comment.CreatedAt,
			&authorID,
			&authorUsername,
			&authorEmail,
			&authorDisplayName,
			&authorAvatarURL,
			&authorCreatedAt,
			&authorUpdatedAt,
		)
		if err != nil {
			return nil, utils.WrapError(err, "failed to scan comment row")
		}

		if authorID.Valid {
			authorUUID, _ := uuid.Parse(authorID.String)
			author.ID = authorUUID
			author.Username = authorUsername.String
			if authorEmail.Valid {
				author.Email = &authorEmail.String
			}
			if authorDisplayName.Valid {
				author.DisplayName = &authorDisplayName.String
			}
			if authorAvatarURL.Valid {
				author.AvatarURL = &authorAvatarURL.String
			}
			author.CreatedAt = authorCreatedAt.Time
			author.UpdatedAt = authorUpdatedAt.Time
			comment.Author = &author
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.WrapError(err, "error iterating comment rows")
	}

	post.Comments = comments
	return post, nil
}

// Update updates a post's information
func (r *postRepository) Update(id uuid.UUID, updates *models.UpdatePostRequest) (*models.Post, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, updates.Title)
		argIndex++
	}

	if updates.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, updates.Content)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByIDWithAuthor(id)
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE posts 
		SET %s
		WHERE id = $%d AND deleted_at IS NULL`,
		strings.Join(setParts, ", "),
		argIndex,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, utils.WrapError(err, "failed to update post")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return nil, utils.ErrPostNotFound
	}

	return r.GetByIDWithAuthor(id)
}

// Delete soft deletes a post
func (r *postRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE posts 
		SET deleted_at = $1 
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return utils.WrapError(err, "failed to delete post")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return utils.ErrPostNotFound
	}

	return nil
}

// List retrieves a paginated list of posts with authors
func (r *postRepository) List(limit, offset int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_by, p.created_at, p.updated_at,
		       u.id, u.username, u.email, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.created_by = u.id
		WHERE p.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list posts")
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var author models.User

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedBy,
			&post.CreatedAt,
			&post.UpdatedAt,
			&author.ID,
			&author.Username,
			&author.Email,
			&author.DisplayName,
			&author.AvatarURL,
			&author.CreatedAt,
			&author.UpdatedAt,
		)
		if err != nil {
			return nil, utils.WrapError(err, "failed to scan post row")
		}

		post.Author = &author
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.WrapError(err, "error iterating post rows")
	}

	return posts, nil
}

// ListByUser retrieves a paginated list of posts by a specific user
func (r *postRepository) ListByUser(userID uuid.UUID, limit, offset int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.created_by, p.created_at, p.updated_at,
		       u.id, u.username, u.email, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.created_by = u.id
		WHERE p.created_by = $1 AND p.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list posts by user")
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var author models.User

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedBy,
			&post.CreatedAt,
			&post.UpdatedAt,
			&author.ID,
			&author.Username,
			&author.Email,
			&author.DisplayName,
			&author.AvatarURL,
			&author.CreatedAt,
			&author.UpdatedAt,
		)
		if err != nil {
			return nil, utils.WrapError(err, "failed to scan post row")
		}

		post.Author = &author
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.WrapError(err, "error iterating post rows")
	}

	return posts, nil
}
