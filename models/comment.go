package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment in the system with support for nested comments
type Comment struct {
	ID        uuid.UUID   `json:"id" db:"id"`
	Content   string      `json:"content" db:"content"`
	PostID    uuid.UUID   `json:"post_id" db:"post_id"`
	ParentID  *uuid.UUID  `json:"parent_id" db:"parent_id"`
	Path      []uuid.UUID `json:"path" db:"path"`
	ThreadID  uuid.UUID   `json:"thread_id" db:"thread_id"`
	CreatedBy *uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	DeletedAt *time.Time  `json:"-" db:"deleted_at"`

	// Associations (loaded separately)
	Post     *Post     `json:"post,omitempty"`
	Parent   *Comment  `json:"parent,omitempty"`
	Author   *User     `json:"author,omitempty"`
	Children []Comment `json:"children,omitempty"`
}

// CreateCommentRequest represents the request payload for creating a comment
type CreateCommentRequest struct {
	Content  string     `json:"content" validate:"required,min=1"`
	PostID   uuid.UUID  `json:"post_id" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id" validate:"omitempty,uuid"`
}

// UpdateCommentRequest represents the request payload for updating a comment
type UpdateCommentRequest struct {
	Content *string `json:"content" validate:"omitempty,min=1"`
}

// GetCommentRequest represents the request payload for getting a comment by ID
type GetCommentRequest struct {
	ID string `json:"id" validate:"required,uuid" uri:"id"`
}

// DeleteCommentRequest represents the request payload for deleting a comment
type DeleteCommentRequest struct {
	ID string `json:"id" validate:"required,uuid" uri:"id"`
}

// ListCommentsRequest represents the request payload for listing comments
type ListCommentsRequest struct {
	PostID string `json:"post_id" validate:"required,uuid" uri:"postId"`
	Limit  int    `json:"limit" validate:"omitempty,gte=1,lte=100" form:"limit"`
	Offset int    `json:"offset" validate:"omitempty,gte=0" form:"offset"`
}

// CommentResponse represents the response payload for comment data
type CommentResponse struct {
	ID        uuid.UUID         `json:"id"`
	Content   string            `json:"content"`
	PostID    uuid.UUID         `json:"post_id"`
	ParentID  *uuid.UUID        `json:"parent_id"`
	Path      []uuid.UUID       `json:"path"`
	ThreadID  uuid.UUID         `json:"thread_id"`
	CreatedBy *uuid.UUID        `json:"created_by"`
	Author    *UserResponse     `json:"author,omitempty"`
	Children  []CommentResponse `json:"children,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// ToResponse converts Comment model to CommentResponse
func (c *Comment) ToResponse() CommentResponse {
	var author *UserResponse
	if c.Author != nil {
		authorResp := c.Author.ToResponse()
		author = &authorResp
	}

	children := make([]CommentResponse, len(c.Children))
	for i, child := range c.Children {
		children[i] = child.ToResponse()
	}

	return CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		PostID:    c.PostID,
		ParentID:  c.ParentID,
		Path:      c.Path,
		ThreadID:  c.ThreadID,
		CreatedBy: c.CreatedBy,
		Author:    author,
		Children:  children,
		CreatedAt: c.CreatedAt,
	}
}
