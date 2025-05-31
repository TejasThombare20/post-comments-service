package services

import (
	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/repository"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/google/uuid"
)

// PostService interface defines post business logic methods
type PostService interface {
	CreatePost(req *models.CreatePostRequest, userID uuid.UUID) (*models.Post, error)
	GetPostByID(id uuid.UUID) (*models.Post, error)
	GetPostWithComments(id uuid.UUID) (*models.Post, error)
	UpdatePost(id uuid.UUID, req *models.UpdatePostRequest, userID uuid.UUID) (*models.Post, error)
	DeletePost(id uuid.UUID, userID uuid.UUID) error
	ListPosts(limit, offset int) ([]models.Post, error)
	ListPostsByUser(userID uuid.UUID, limit, offset int) ([]models.Post, error)
}

// postService implements PostService interface
type postService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

// NewPostService creates a new post service instance
func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) PostService {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// CreatePost creates a new post
func (s *postService) CreatePost(req *models.CreatePostRequest, userID uuid.UUID) (*models.Post, error) {
	// Verify user exists
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, err
	}

	// Create post model
	post := &models.Post{
		ID:        uuid.New(),
		Title:     req.Title,
		Content:   req.Content,
		CreatedBy: userID,
	}

	// Save post to database
	if err := s.postRepo.Create(post); err != nil {
		return nil, utils.WrapError(err, "failed to create post")
	}

	// Return post with author information
	createdPost, err := s.postRepo.GetByIDWithAuthor(post.ID)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get created post with author")
	}

	return createdPost, nil
}

// GetPostByID retrieves a post by ID with author information
func (s *postService) GetPostByID(id uuid.UUID) (*models.Post, error) {
	post, err := s.postRepo.GetByIDWithAuthor(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// GetPostWithComments retrieves a post by ID with comments and authors
func (s *postService) GetPostWithComments(id uuid.UUID) (*models.Post, error) {
	post, err := s.postRepo.GetByIDWithComments(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// UpdatePost updates a post (only by the author)
func (s *postService) UpdatePost(id uuid.UUID, req *models.UpdatePostRequest, userID uuid.UUID) (*models.Post, error) {
	// Get existing post
	existingPost, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user is the author
	if existingPost.CreatedBy != userID {
		return nil, utils.ErrForbidden
	}

	// Update post
	updatedPost, err := s.postRepo.Update(id, req)
	if err != nil {
		return nil, utils.WrapError(err, "failed to update post")
	}

	return updatedPost, nil
}

// DeletePost deletes a post (only by the author)
func (s *postService) DeletePost(id uuid.UUID, userID uuid.UUID) error {
	// Get existing post
	existingPost, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if user is the author
	if existingPost.CreatedBy != userID {
		return utils.ErrForbidden
	}

	// Delete post
	if err := s.postRepo.Delete(id); err != nil {
		return utils.WrapError(err, "failed to delete post")
	}

	return nil
}

// ListPosts retrieves a paginated list of posts with authors
func (s *postService) ListPosts(limit, offset int) ([]models.Post, error) {
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

	posts, err := s.postRepo.List(limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list posts")
	}

	return posts, nil
}

// ListPostsByUser retrieves a paginated list of posts by a specific user
func (s *postService) ListPostsByUser(userID uuid.UUID, limit, offset int) ([]models.Post, error) {
	// Verify user exists
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, err
	}

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

	posts, err := s.postRepo.ListByUser(userID, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list posts by user")
	}

	return posts, nil
}
