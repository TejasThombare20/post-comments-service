package models

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a post in the system
type Post struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Content   string     `json:"content" db:"content"`
	CreatedBy uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`

	// Associations (loaded separately)
	Author   *User     `json:"author,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// CreatePostRequest represents the request payload for creating a post
type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// UpdatePostRequest represents the request payload for updating a post
type UpdatePostRequest struct {
	Title   *string `json:"title" validate:"omitempty,min=1,max=200"`
	Content *string `json:"content" validate:"omitempty,min=1"`
}

// GetPostRequest represents the request payload for getting a post by ID
type GetPostRequest struct {
	ID string `json:"id" validate:"required,uuid" uri:"id"`
}

// DeletePostRequest represents the request payload for deleting a post
type DeletePostRequest struct {
	ID string `json:"id" validate:"required,uuid" uri:"id"`
}

// ListPostsRequest represents the request payload for listing posts
type ListPostsRequest struct {
	Limit  int `json:"limit" validate:"omitempty,gte=1,lte=100" form:"limit"`
	Offset int `json:"offset" validate:"omitempty,gte=0" form:"offset"`
}

// ListPostsByUserRequest represents the request payload for listing posts by user
type ListPostsByUserRequest struct {
	UserID string `json:"user_id" validate:"required,uuid" uri:"userId"`
	Limit  int    `json:"limit" validate:"omitempty,gte=1,lte=100" form:"limit"`
	Offset int    `json:"offset" validate:"omitempty,gte=0" form:"offset"`
}

// PostResponse represents the response payload for post data
type PostResponse struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	CreatedBy uuid.UUID    `json:"created_by"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// PostWithCommentsResponse represents the response payload for post data with comments
type PostWithCommentsResponse struct {
	ID        uuid.UUID         `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	CreatedBy uuid.UUID         `json:"created_by"`
	Author    UserResponse      `json:"author"`
	Comments  []CommentResponse `json:"comments"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// ToResponse converts Post model to PostResponse
func (p *Post) ToResponse() PostResponse {
	var author UserResponse
	if p.Author != nil {
		author = p.Author.ToResponse()
	}

	return PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedBy: p.CreatedBy,
		Author:    author,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToResponseWithComments converts Post model to PostWithCommentsResponse
func (p *Post) ToResponseWithComments() PostWithCommentsResponse {
	var author UserResponse
	if p.Author != nil {
		author = p.Author.ToResponse()
	}

	comments := make([]CommentResponse, len(p.Comments))
	for i, comment := range p.Comments {
		comments[i] = comment.ToResponse()
	}

	return PostWithCommentsResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedBy: p.CreatedBy,
		Author:    author,
		Comments:  comments,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
