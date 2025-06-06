package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// CommentRepository interface defines comment data access methods
type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByID(id uuid.UUID) (*models.Comment, error)
	GetByIDWithAuthor(id uuid.UUID) (*models.Comment, error)
	Update(id uuid.UUID, updates *models.UpdateCommentRequest) (*models.Comment, error)
	Delete(id uuid.UUID) error
	ListByPost(postID uuid.UUID, limit, offset int) ([]models.Comment, error)
	GetReplies(parentID uuid.UUID, limit, offset int) ([]models.Comment, error)
	IncrementRepliesCount(commentID uuid.UUID) error
}

// commentRepository implements CommentRepository interface
type commentRepository struct {
	db *sql.DB
}

// NewCommentRepository creates a new comment repository instance
func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

// convertUUIDSliceToStringArray converts []uuid.UUID to pq.StringArray
func convertUUIDSliceToStringArray(uuids []uuid.UUID) pq.StringArray {
	strings := make([]string, len(uuids))
	for i, u := range uuids {
		strings[i] = u.String()
	}
	return pq.StringArray(strings)
}

// convertStringArrayToUUIDSlice converts pq.StringArray to []uuid.UUID
func convertStringArrayToUUIDSlice(strings pq.StringArray) []uuid.UUID {
	uuids := make([]uuid.UUID, len(strings))
	for i, s := range strings {
		if u, err := uuid.Parse(s); err == nil {
			uuids[i] = u
		}
	}
	return uuids
}

// Create creates a new comment in the database
func (r *commentRepository) Create(comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, content, post_id, parent_id, path, thread_id, created_by, created_at, updated_at, replies_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	pathArray := convertUUIDSliceToStringArray(comment.Path)

	_, err := r.db.Exec(query,
		comment.ID,
		comment.Content,
		comment.PostID,
		comment.ParentID,
		pathArray,
		comment.ThreadID,
		comment.CreatedBy,
		comment.CreatedAt,
		comment.UpdatedAt,
		comment.RepliesCount,
	)

	if err != nil {
		return utils.WrapError(err, "failed to create comment")
	}

	return nil
}

// GetByID retrieves a comment by ID
func (r *commentRepository) GetByID(id uuid.UUID) (*models.Comment, error) {
	query := `
		SELECT id, content, post_id, parent_id, path, thread_id, created_by, created_at, updated_at, replies_count
		FROM comments 
		WHERE id = $1 AND deleted_at IS NULL`

	var comment models.Comment
	var pathArray pq.StringArray

	err := r.db.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.PostID,
		&comment.ParentID,
		&pathArray,
		&comment.ThreadID,
		&comment.CreatedBy,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.RepliesCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrCommentNotFound
		}
		return nil, utils.WrapError(err, "failed to get comment by ID")
	}

	comment.Path = convertStringArrayToUUIDSlice(pathArray)
	return &comment, nil
}

// GetByIDWithAuthor retrieves a comment by ID with author information
func (r *commentRepository) GetByIDWithAuthor(id uuid.UUID) (*models.Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.parent_id, c.path, c.thread_id, c.created_by, c.created_at, c.updated_at, c.replies_count,
		       u.id, u.username, u.email, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM comments c
		JOIN users u ON c.created_by = u.id
		WHERE c.id = $1 AND c.deleted_at IS NULL AND u.deleted_at IS NULL`

	var comment models.Comment
	var author models.User
	var pathArray pq.StringArray

	err := r.db.QueryRow(query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.PostID,
		&comment.ParentID,
		&pathArray,
		&comment.ThreadID,
		&comment.CreatedBy,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.RepliesCount,
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
			return nil, utils.ErrCommentNotFound
		}
		return nil, utils.WrapError(err, "failed to get comment with author")
	}

	comment.Path = convertStringArrayToUUIDSlice(pathArray)
	comment.Author = &author
	return &comment, nil
}

// Update updates a comment's information
func (r *commentRepository) Update(id uuid.UUID, updates *models.UpdateCommentRequest) (*models.Comment, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, updates.Content)
		argIndex++
	}

	// Always update updated_at when any field is updated
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	if len(setParts) == 1 { // Only updated_at was added
		// No actual content updates to perform, just return the current comment
		return r.GetByIDWithAuthor(id)
	}

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE comments 
		SET %s
		WHERE id = $%d AND deleted_at IS NULL`,
		strings.Join(setParts, ", "),
		argIndex,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, utils.WrapError(err, "failed to update comment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return nil, utils.ErrCommentNotFound
	}

	// Return updated comment with author
	return r.GetByIDWithAuthor(id)
}

// Delete soft deletes a comment
func (r *commentRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE comments 
		SET deleted_at = $1 
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return utils.WrapError(err, "failed to delete comment")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return utils.ErrCommentNotFound
	}

	return nil
}

// ListByPost retrieves a paginated list of comments for a specific post
func (r *commentRepository) ListByPost(postID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.parent_id, c.path, c.thread_id, c.created_by, c.created_at, c.updated_at, c.replies_count,
		       u.id, u.username, u.email, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM comments c
		LEFT JOIN users u ON c.created_by = u.id AND u.deleted_at IS NULL
		WHERE c.post_id = $1 AND c.deleted_at IS NULL AND c.parent_id IS NULL
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list comments by post")
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var author models.User
		var pathArray pq.StringArray
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
			&pathArray,
			&comment.ThreadID,
			&comment.CreatedBy,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.RepliesCount,
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

		// Convert PostgreSQL string array to UUID slice
		comment.Path = convertStringArrayToUUIDSlice(pathArray)

		// Set author if exists
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

	return comments, nil
}

// GetReplies retrieves replies to a specific comment
func (r *commentRepository) GetReplies(parentID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.parent_id, c.path, c.thread_id, c.created_by, c.created_at, c.updated_at, c.replies_count,
		       u.id, u.username, u.email, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM comments c
		LEFT JOIN users u ON c.created_by = u.id AND u.deleted_at IS NULL
		WHERE c.parent_id = $1 AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, parentID, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get comment replies")
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var author models.User
		var pathArray pq.StringArray
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
			&pathArray,
			&comment.ThreadID,
			&comment.CreatedBy,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.RepliesCount,
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

		// Convert PostgreSQL string array to UUID slice
		comment.Path = convertStringArrayToUUIDSlice(pathArray)

		// Set author if exists
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

	return comments, nil
}

// IncrementRepliesCount increments the replies count for a comment
// NOTE: This should only be used for manual corrections or data migration.
// Normal reply creation automatically increments the count via database trigger.
func (r *commentRepository) IncrementRepliesCount(commentID uuid.UUID) error {
	query := `
		UPDATE comments 
		SET replies_count = replies_count + 1, updated_at = $1 
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, time.Now(), commentID)
	if err != nil {
		return utils.WrapError(err, "failed to increment replies count")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.WrapError(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return utils.ErrCommentNotFound
	}

	return nil
}
