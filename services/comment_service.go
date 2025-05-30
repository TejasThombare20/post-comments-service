package services

import (
	"time"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/repository"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/TejasThombare20/post-comments-service/validator"
	"github.com/google/uuid"
)

// CommentService interface defines comment business logic methods
type CommentService interface {
	CreateComment(userID uuid.UUID, req *models.CreateCommentRequest) (*models.Comment, error)
	GetCommentByID(req *models.GetCommentRequest) (*models.Comment, error)
	UpdateComment(id uuid.UUID, userID uuid.UUID, req *models.UpdateCommentRequest) (*models.Comment, error)
	DeleteComment(req *models.DeleteCommentRequest, userID uuid.UUID) error
	ListCommentsByPost(req *models.ListCommentsRequest) ([]models.Comment, error)
	GetCommentReplies(commentID uuid.UUID, limit, offset int) ([]models.Comment, error)
}

// commentService implements CommentService interface
type commentService struct {
	commentRepo   repository.CommentRepository
	postRepo      repository.PostRepository
	userRepo      repository.UserRepository
	validator     *validator.Validator
	htmlSanitizer *utils.HTMLSanitizer
}

// NewCommentService creates a new comment service instance
func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRepo repository.UserRepository, validator *validator.Validator) CommentService {
	return &commentService{
		commentRepo:   commentRepo,
		postRepo:      postRepo,
		userRepo:      userRepo,
		validator:     validator,
		htmlSanitizer: utils.NewHTMLSanitizer(),
	}
}

// CreateComment creates a new comment or reply
func (s *commentService) CreateComment(userID uuid.UUID, req *models.CreateCommentRequest) (*models.Comment, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	postID, err := uuid.Parse(*req.PostID)
	if err != nil {
		return nil, utils.WrapError(err, "invalid post_id format")
	}

	if _, err = s.userRepo.GetByID(userID); err != nil {
		return nil, utils.WrapError(err, "failed to find user")
	}

	if _, err = s.postRepo.GetByID(postID); err != nil {
		return nil, utils.WrapError(err, "failed to find post")
	}

	if err := s.htmlSanitizer.ValidateHTMLContent(*req.Content); err != nil {
		return nil, utils.WrapError(utils.ErrInvalidInput, "invalid HTML content: "+err.Error())
	}

	sanitizedContent := s.htmlSanitizer.ProcessCommentContent(*req.Content)

	comment := &models.Comment{
		ID:        uuid.New(),
		Content:   sanitizedContent,
		PostID:    postID,
		CreatedBy: &userID,
		CreatedAt: time.Now(),
		Path:      []uuid.UUID{},
	}

	if req.ParentID != nil && *req.ParentID != "" {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, utils.WrapError(err, "invalid parent_id format")
		}

		parentComment, err := s.commentRepo.GetByID(parentID)
		if err != nil {
			return nil, utils.WrapError(err, "failed to find parent comment")
		}

		if parentComment.PostID != postID {
			return nil, utils.WrapError(utils.ErrInvalidInput, "parent comment does not belong to the same post")
		}

		comment.ParentID = &parentID
		comment.ThreadID = parentComment.ThreadID
		comment.Path = append(parentComment.Path, parentID)
	} else {
		comment.ThreadID = comment.ID
		comment.Path = append(comment.Path, comment.ID)
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, utils.WrapError(err, "failed to create comment")
	}

	createdComment, err := s.commentRepo.GetByIDWithAuthor(comment.ID)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get created comment with author")
	}

	return createdComment, nil
}

// GetCommentByID retrieves a comment by ID with author information
func (s *commentService) GetCommentByID(req *models.GetCommentRequest) (*models.Comment, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	commentID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, utils.WrapError(err, "invalid comment ID format")
	}

	comment, err := s.commentRepo.GetByIDWithAuthor(commentID)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get comment by ID")
	}

	return comment, nil
}

// UpdateComment updates a comment's content
func (s *commentService) UpdateComment(id uuid.UUID, userID uuid.UUID, req *models.UpdateCommentRequest) (*models.Comment, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	existingComment, err := s.commentRepo.GetByID(id)
	if err != nil {
		return nil, utils.WrapError(err, "failed to find comment for update")
	}

	if existingComment.CreatedBy == nil || *existingComment.CreatedBy != userID {
		return nil, utils.ErrForbidden
	}

	if req.Content != nil {
		if err := s.htmlSanitizer.ValidateHTMLContent(*req.Content); err != nil {
			return nil, utils.WrapError(utils.ErrInvalidInput, "invalid HTML content: "+err.Error())
		}
		sanitizedContent := s.htmlSanitizer.ProcessCommentContent(*req.Content)
		req.Content = &sanitizedContent
	}

	updatedComment, err := s.commentRepo.Update(id, req)
	if err != nil {
		return nil, utils.WrapError(err, "failed to update comment")
	}

	return updatedComment, nil
}

// DeleteComment deletes a comment
func (s *commentService) DeleteComment(req *models.DeleteCommentRequest, userID uuid.UUID) error {
	if err := s.validator.ValidateStruct(req); err != nil {
		return err
	}

	commentID, err := uuid.Parse(req.ID)
	if err != nil {
		return utils.WrapError(err, "invalid comment ID format")
	}

	existingComment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return utils.WrapError(err, "failed to find comment for deletion")
	}

	if existingComment.CreatedBy == nil || *existingComment.CreatedBy != userID {
		return utils.ErrForbidden
	}

	if err := s.commentRepo.Delete(commentID); err != nil {
		return utils.WrapError(err, "failed to delete comment")
	}

	return nil
}

// ListCommentsByPost retrieves comments for a specific post
func (s *commentService) ListCommentsByPost(req *models.ListCommentsRequest) ([]models.Comment, error) {
	if err := s.validator.ValidateStruct(req); err != nil {
		return nil, err
	}

	postID, err := uuid.Parse(req.PostID)
	if err != nil {
		return nil, utils.WrapError(err, "invalid post ID format")
	}

	if _, err = s.postRepo.GetByID(postID); err != nil {
		return nil, utils.WrapError(err, "failed to find post")
	}

	limit := 20
	offset := 0

	if req.Limit > 0 {
		limit = req.Limit
	}
	if req.Offset >= 0 {
		offset = req.Offset
	}
	if limit > 100 {
		limit = 100
	}

	comments, err := s.commentRepo.ListByPost(postID, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to list comments by post")
	}

	return comments, nil
}

// GetCommentReplies retrieves replies for a specific comment
func (s *commentService) GetCommentReplies(commentID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	if _, err := s.commentRepo.GetByID(commentID); err != nil {
		return nil, utils.WrapError(err, "failed to find comment")
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if limit > 100 {
		limit = 100
	}

	replies, err := s.commentRepo.GetReplies(commentID, limit, offset)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get comment replies")
	}

	return replies, nil
}
