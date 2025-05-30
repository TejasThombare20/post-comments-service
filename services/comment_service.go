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
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	validator   *validator.Validator
}

// NewCommentService creates a new comment service instance
func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRepo repository.UserRepository, validator *validator.Validator) CommentService {
	utils.LogInfo("Initializing comment service", utils.LogFields{
		"component": "comment_service",
	})
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		validator:   validator,
	}
}

// CreateComment creates a new comment or reply
func (s *commentService) CreateComment(userID uuid.UUID, req *models.CreateCommentRequest) (*models.Comment, error) {
	utils.LogInfo("Creating comment", utils.LogFields{
		"user_id":  userID,
		"post_id":  req.PostID,
		"is_reply": req.ParentID != nil,
	})

	// Validate request
	if err := s.validator.ValidateStruct(req); err != nil {
		utils.LogError("Comment creation validation failed", err, utils.LogFields{
			"user_id": userID,
			"post_id": req.PostID,
		})
		return nil, err
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		utils.LogError("User not found for comment creation", err, utils.LogFields{
			"user_id": userID,
		})
		return nil, utils.WrapError(err, "failed to find user")
	}

	// Check if post exists
	_, err = s.postRepo.GetByID(req.PostID)
	if err != nil {
		utils.LogError("Post not found for comment creation", err, utils.LogFields{
			"post_id": req.PostID,
		})
		return nil, utils.WrapError(err, "failed to find post")
	}

	// Create comment model
	comment := &models.Comment{
		ID:        uuid.New(),
		Content:   req.Content,
		PostID:    req.PostID,
		CreatedBy: &userID,
		CreatedAt: time.Now(),
		Path:      []uuid.UUID{},
	}

	// Handle parent comment if this is a reply
	if req.ParentID != nil {
		utils.LogInfo("Processing reply comment", utils.LogFields{
			"comment_id": comment.ID,
			"parent_id":  *req.ParentID,
		})

		// Get parent comment to build path
		parentComment, err := s.commentRepo.GetByID(*req.ParentID)
		if err != nil {
			utils.LogError("Parent comment not found", err, utils.LogFields{
				"parent_id": *req.ParentID,
			})
			return nil, utils.WrapError(err, "failed to find parent comment")
		}

		// Verify parent comment belongs to the same post
		if parentComment.PostID != req.PostID {
			utils.LogError("Parent comment belongs to different post", nil, utils.LogFields{
				"parent_id":         *req.ParentID,
				"parent_post_id":    parentComment.PostID,
				"requested_post_id": req.PostID,
			})
			return nil, utils.WrapError(utils.ErrInvalidInput, "parent comment does not belong to the same post")
		}

		comment.ParentID = req.ParentID
		comment.ThreadID = parentComment.ThreadID

		// Build path: parent's path + parent's ID
		comment.Path = append(parentComment.Path, *req.ParentID)
	} else {
		// This is a top-level comment, it becomes its own thread
		comment.ThreadID = comment.ID
		utils.LogInfo("Creating top-level comment", utils.LogFields{
			"comment_id": comment.ID,
			"thread_id":  comment.ThreadID,
		})
	}

	// Save comment to database
	if err := s.commentRepo.Create(comment); err != nil {
		utils.LogError("Failed to save comment to database", err, utils.LogFields{
			"comment_id": comment.ID,
			"user_id":    userID,
		})
		return nil, utils.WrapError(err, "failed to create comment")
	}

	// Get comment with author information
	createdComment, err := s.commentRepo.GetByIDWithAuthor(comment.ID)
	if err != nil {
		utils.LogError("Failed to retrieve created comment with author", err, utils.LogFields{
			"comment_id": comment.ID,
		})
		return nil, utils.WrapError(err, "failed to get created comment with author")
	}

	utils.LogInfo("Comment created successfully", utils.LogFields{
		"comment_id": comment.ID,
		"post_id":    req.PostID,
		"created_by": userID,
		"is_reply":   req.ParentID != nil,
		"thread_id":  comment.ThreadID,
	})

	return createdComment, nil
}

// GetCommentByID retrieves a comment by ID with author information
func (s *commentService) GetCommentByID(req *models.GetCommentRequest) (*models.Comment, error) {
	utils.LogInfo("Retrieving comment by ID", utils.LogFields{
		"comment_id": req.ID,
	})

	// Validate request
	if err := s.validator.ValidateStruct(req); err != nil {
		utils.LogError("Get comment validation failed", err, utils.LogFields{
			"comment_id": req.ID,
		})
		return nil, err
	}

	// Parse UUID from string
	commentID, err := uuid.Parse(req.ID)
	if err != nil {
		utils.LogError("Invalid comment ID format", err, utils.LogFields{
			"comment_id": req.ID,
		})
		return nil, utils.WrapError(err, "invalid comment ID format")
	}

	comment, err := s.commentRepo.GetByIDWithAuthor(commentID)
	if err != nil {
		utils.LogError("Failed to retrieve comment", err, utils.LogFields{
			"comment_id": commentID,
		})
		return nil, utils.WrapError(err, "failed to get comment by ID")
	}

	utils.LogInfo("Comment retrieved successfully", utils.LogFields{
		"comment_id": commentID,
		"author_id":  comment.CreatedBy,
	})

	return comment, nil
}

// UpdateComment updates a comment's content
func (s *commentService) UpdateComment(id uuid.UUID, userID uuid.UUID, req *models.UpdateCommentRequest) (*models.Comment, error) {
	utils.LogInfo("Updating comment", utils.LogFields{
		"comment_id": id,
		"user_id":    userID,
	})

	// Validate request
	if err := s.validator.ValidateStruct(req); err != nil {
		utils.LogError("Update comment validation failed", err, utils.LogFields{
			"comment_id": id,
			"user_id":    userID,
		})
		return nil, err
	}

	// Check if comment exists and user owns it
	existingComment, err := s.commentRepo.GetByID(id)
	if err != nil {
		utils.LogError("Comment not found for update", err, utils.LogFields{
			"comment_id": id,
		})
		return nil, utils.WrapError(err, "failed to find comment for update")
	}

	// Check if user owns the comment
	if existingComment.CreatedBy == nil || *existingComment.CreatedBy != userID {
		utils.LogError("User not authorized to update comment", nil, utils.LogFields{
			"comment_id":    id,
			"user_id":       userID,
			"comment_owner": existingComment.CreatedBy,
		})
		return nil, utils.ErrForbidden
	}

	// Update comment
	updatedComment, err := s.commentRepo.Update(id, req)
	if err != nil {
		utils.LogError("Failed to update comment in database", err, utils.LogFields{
			"comment_id": id,
			"user_id":    userID,
		})
		return nil, utils.WrapError(err, "failed to update comment")
	}

	utils.LogInfo("Comment updated successfully", utils.LogFields{
		"comment_id": updatedComment.ID,
		"updated_by": userID,
	})

	return updatedComment, nil
}

// DeleteComment deletes a comment
func (s *commentService) DeleteComment(req *models.DeleteCommentRequest, userID uuid.UUID) error {
	utils.LogInfo("Deleting comment", utils.LogFields{
		"comment_id": req.ID,
		"user_id":    userID,
	})

	// Validate request
	if err := s.validator.ValidateStruct(req); err != nil {
		utils.LogError("Delete comment validation failed", err, utils.LogFields{
			"comment_id": req.ID,
			"user_id":    userID,
		})
		return err
	}

	// Parse UUID from string
	commentID, err := uuid.Parse(req.ID)
	if err != nil {
		utils.LogError("Invalid comment ID format for deletion", err, utils.LogFields{
			"comment_id": req.ID,
		})
		return utils.WrapError(err, "invalid comment ID format")
	}

	// Check if comment exists and user owns it
	existingComment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		utils.LogError("Comment not found for deletion", err, utils.LogFields{
			"comment_id": commentID,
		})
		return utils.WrapError(err, "failed to find comment for deletion")
	}

	// Check if user owns the comment
	if existingComment.CreatedBy == nil || *existingComment.CreatedBy != userID {
		utils.LogError("User not authorized to delete comment", nil, utils.LogFields{
			"comment_id":    commentID,
			"user_id":       userID,
			"comment_owner": existingComment.CreatedBy,
		})
		return utils.ErrForbidden
	}

	// Delete comment
	if err := s.commentRepo.Delete(commentID); err != nil {
		utils.LogError("Failed to delete comment from database", err, utils.LogFields{
			"comment_id": commentID,
			"user_id":    userID,
		})
		return utils.WrapError(err, "failed to delete comment")
	}

	utils.LogInfo("Comment deleted successfully", utils.LogFields{
		"comment_id": commentID,
		"deleted_by": userID,
	})

	return nil
}

// ListCommentsByPost retrieves comments for a specific post
func (s *commentService) ListCommentsByPost(req *models.ListCommentsRequest) ([]models.Comment, error) {
	utils.LogInfo("Listing comments by post", utils.LogFields{
		"post_id": req.PostID,
		"limit":   req.Limit,
		"offset":  req.Offset,
	})

	// Validate request
	if err := s.validator.ValidateStruct(req); err != nil {
		utils.LogError("List comments validation failed", err, utils.LogFields{
			"post_id": req.PostID,
		})
		return nil, err
	}

	// Parse post ID
	postID, err := uuid.Parse(req.PostID)
	if err != nil {
		utils.LogError("Invalid post ID format", err, utils.LogFields{
			"post_id": req.PostID,
		})
		return nil, utils.WrapError(err, "invalid post ID format")
	}

	// Check if post exists
	_, err = s.postRepo.GetByID(postID)
	if err != nil {
		utils.LogError("Post not found for comment listing", err, utils.LogFields{
			"post_id": postID,
		})
		return nil, utils.WrapError(err, "failed to find post")
	}

	// Set default values if not provided
	limit := 20
	offset := 0

	if req.Limit > 0 {
		limit = req.Limit
	}
	if req.Offset >= 0 {
		offset = req.Offset
	}

	// Enforce maximum limit
	if limit > 100 {
		limit = 100
		utils.LogInfo("Comment list limit capped at maximum", utils.LogFields{
			"requested_limit": req.Limit,
			"applied_limit":   limit,
		})
	}

	comments, err := s.commentRepo.ListByPost(postID, limit, offset)
	if err != nil {
		utils.LogError("Failed to retrieve comments from database", err, utils.LogFields{
			"post_id": postID,
			"limit":   limit,
			"offset":  offset,
		})
		return nil, utils.WrapError(err, "failed to list comments by post")
	}

	utils.LogInfo("Comments retrieved successfully", utils.LogFields{
		"post_id": postID,
		"count":   len(comments),
		"limit":   limit,
		"offset":  offset,
	})

	return comments, nil
}

// GetCommentReplies retrieves replies for a specific comment
func (s *commentService) GetCommentReplies(commentID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	utils.LogInfo("Getting comment replies", utils.LogFields{
		"comment_id": commentID,
		"limit":      limit,
		"offset":     offset,
	})

	// Check if comment exists
	_, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		utils.LogError("Parent comment not found for replies", err, utils.LogFields{
			"comment_id": commentID,
		})
		return nil, utils.WrapError(err, "failed to find comment")
	}

	// Set default values if not provided
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Enforce maximum limit
	if limit > 100 {
		limit = 100
		utils.LogInfo("Comment replies limit capped at maximum", utils.LogFields{
			"requested_limit": limit,
			"applied_limit":   100,
		})
	}

	replies, err := s.commentRepo.GetReplies(commentID, limit, offset)
	if err != nil {
		utils.LogError("Failed to retrieve comment replies", err, utils.LogFields{
			"comment_id": commentID,
			"limit":      limit,
			"offset":     offset,
		})
		return nil, utils.WrapError(err, "failed to get comment replies")
	}

	utils.LogInfo("Comment replies retrieved successfully", utils.LogFields{
		"comment_id":    commentID,
		"replies_count": len(replies),
		"limit":         limit,
		"offset":        offset,
	})

	return replies, nil
}
