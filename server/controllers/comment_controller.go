package controllers

import (
	"net/http"
	"strconv"

	"github.com/TejasThombare20/post-comments-service/models"
	"github.com/TejasThombare20/post-comments-service/services"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CommentController handles comment-related HTTP requests
type CommentController struct {
	commentService services.CommentService
}

// NewCommentController creates a new comment controller instance
func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

// CreateComment handles POST /posts/:postId/comments
func (cc *CommentController) CreateComment(c *gin.Context) {
	postIDParam := c.Param("postId")
	_, err := uuid.Parse(postIDParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid post ID format")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	req.PostID = &postIDParam

	comment, err := cc.commentService.CreateComment(userID, &req)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "Post")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, comment)
}

// GetComment handles GET /comments/:id
func (cc *CommentController) GetComment(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid comment ID format")
		return
	}

	req := &models.GetCommentRequest{ID: idParam}
	comment, err := cc.commentService.GetCommentByID(req)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "Comment")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, comment)
}

// UpdateComment handles PUT /comments/:id
func (cc *CommentController) UpdateComment(c *gin.Context) {
	idParam := c.Param("id")
	commentID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid comment ID format")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	comment, err := cc.commentService.UpdateComment(commentID, userID, &req)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "Comment")
			return
		}
		if utils.IsForbiddenError(err) {
			utils.ForbiddenResponse(c, "You can only update your own comments")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, comment)
}

// DeleteComment handles DELETE /comments/:id
func (cc *CommentController) DeleteComment(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid comment ID format")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	req := &models.DeleteCommentRequest{ID: idParam}
	err = cc.commentService.DeleteComment(req, userID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "Comment")
			return
		}
		if utils.IsForbiddenError(err) {
			utils.ForbiddenResponse(c, "You can only delete your own comments")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// ListCommentsByPost handles GET /posts/:postId/comments
func (cc *CommentController) ListCommentsByPost(c *gin.Context) {
	postIDParam := c.Param("postId")
	if postIDParam == "" {
		utils.ValidationErrorResponse(c, "Post ID is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		utils.ValidationErrorResponse(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.ValidationErrorResponse(c, "Invalid offset parameter")
		return
	}

	req := &models.ListCommentsRequest{
		PostID: postIDParam,
		Limit:  limit,
		Offset: offset,
	}

	comments, err := cc.commentService.ListCommentsByPost(req)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "Post")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	commentResponses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = comment.ToResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"comments": commentResponses,
		"limit":    limit,
		"offset":   offset,
		"count":    len(commentResponses),
	})
}

// GetCommentReplies handles GET /comments/:id/replies
func (cc *CommentController) GetCommentReplies(c *gin.Context) {
	idParam := c.Param("id")
	commentID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid comment ID format")
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		utils.ValidationErrorResponse(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.ValidationErrorResponse(c, "Invalid offset parameter")
		return
	}

	replies, err := cc.commentService.GetCommentReplies(commentID, limit, offset)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "Comment")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	replyResponses := make([]models.CommentResponse, len(replies))
	for i, reply := range replies {
		replyResponses[i] = reply.ToResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"replies": replyResponses,
		"limit":   limit,
		"offset":  offset,
		"count":   len(replyResponses),
	})
}

// GetCommentsByPost handles GET /posts/:postId/comments
func (cc *CommentController) GetCommentsByPost(c *gin.Context) {
	postIDParam := c.Param("postId")
	utils.LogInfo("Getting comments by post", utils.LogFields{
		"post_id": postIDParam,
	})

	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		utils.LogError("Invalid post ID format for comments query", err, utils.LogFields{
			"post_id": postIDParam,
		})
		utils.ValidationErrorResponse(c, "Invalid post ID format")
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		utils.LogError("Invalid limit parameter for comments", err, utils.LogFields{
			"post_id": postID,
			"limit":   limitStr,
		})
		utils.ValidationErrorResponse(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.LogError("Invalid offset parameter for comments", err, utils.LogFields{
			"post_id": postID,
			"offset":  offsetStr,
		})
		utils.ValidationErrorResponse(c, "Invalid offset parameter")
		return
	}

	req := &models.ListCommentsRequest{
		PostID: postIDParam,
		Limit:  limit,
		Offset: offset,
	}

	comments, err := cc.commentService.ListCommentsByPost(req)
	if err != nil {
		utils.LogError("Failed to get comments by post", err, utils.LogFields{
			"post_id": postID,
			"limit":   limit,
			"offset":  offset,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("Comments retrieved successfully", utils.LogFields{
		"post_id": postID,
		"count":   len(comments),
		"limit":   limit,
		"offset":  offset,
	})

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"comments": comments,
		"limit":    limit,
		"offset":   offset,
		"count":    len(comments),
	})
}
