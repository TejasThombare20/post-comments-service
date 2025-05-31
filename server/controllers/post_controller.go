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

// PostController handles post-related HTTP requests
type PostController struct {
	postService services.PostService
}

// NewPostController creates a new post controller instance
func NewPostController(postService services.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// CreatePost handles POST /posts
func (pc *PostController) CreatePost(c *gin.Context) {
	utils.LogInfo("Creating new post", utils.LogFields{})

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogError("User not authenticated for post creation", err, utils.LogFields{})
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid request payload for post creation", err, utils.LogFields{
			"user_id": userID,
		})
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	post, err := pc.postService.CreatePost(&req, userID)
	if err != nil {
		utils.LogError("Failed to create post", err, utils.LogFields{
			"user_id": userID,
			"title":   req.Title,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("Post created successfully", utils.LogFields{
		"post_id": post.ID,
		"user_id": userID,
		"title":   post.Title,
	})

	utils.SuccessResponse(c, http.StatusCreated, post)
}

// GetPost handles GET /posts/:id
func (pc *PostController) GetPost(c *gin.Context) {
	idParam := c.Param("id")
	utils.LogInfo("Getting post by ID", utils.LogFields{
		"post_id": idParam,
	})

	postID, err := uuid.Parse(idParam)
	if err != nil {
		utils.LogError("Invalid post ID format", err, utils.LogFields{
			"post_id": idParam,
		})
		utils.ValidationErrorResponse(c, "Invalid post ID format")
		return
	}

	post, err := pc.postService.GetPostByID(postID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("Post not found", err, utils.LogFields{
				"post_id": postID,
			})
			utils.NotFoundResponse(c, "Post")
			return
		}
		utils.LogError("Failed to get post", err, utils.LogFields{
			"post_id": postID,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("Post retrieved successfully", utils.LogFields{
		"post_id": postID,
		"title":   post.Title,
		"author":  post.Author.Username,
	})

	utils.SuccessResponse(c, http.StatusOK, post)
}

// GetPostWithComments handles GET /posts/:id/comments
func (pc *PostController) GetPostWithComments(c *gin.Context) {
	idParam := c.Param("id")
	postID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid post ID format")
		return
	}

	post, err := pc.postService.GetPostWithComments(postID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.NotFoundResponse(c, "Post")
			return
		}
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, post.ToResponseWithComments())
}

// UpdatePost handles PUT /posts/:id
func (pc *PostController) UpdatePost(c *gin.Context) {
	idParam := c.Param("id")
	utils.LogInfo("Updating post", utils.LogFields{
		"post_id": idParam,
	})

	postID, err := uuid.Parse(idParam)
	if err != nil {
		utils.LogError("Invalid post ID format for update", err, utils.LogFields{
			"post_id": idParam,
		})
		utils.ValidationErrorResponse(c, "Invalid post ID format")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogError("User not authenticated for post update", err, utils.LogFields{
			"post_id": postID,
		})
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Invalid request payload for post update", err, utils.LogFields{
			"post_id": postID,
			"user_id": userID,
		})
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	post, err := pc.postService.UpdatePost(postID, &req, userID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("Post not found for update", err, utils.LogFields{
				"post_id": postID,
				"user_id": userID,
			})
			utils.NotFoundResponse(c, "Post")
			return
		}
		if utils.IsForbiddenError(err) {
			utils.LogError("User not authorized to update post", err, utils.LogFields{
				"post_id": postID,
				"user_id": userID,
			})
			utils.ForbiddenResponse(c, "You can only update your own posts")
			return
		}
		utils.LogError("Failed to update post", err, utils.LogFields{
			"post_id": postID,
			"user_id": userID,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("Post updated successfully", utils.LogFields{
		"post_id": postID,
		"user_id": userID,
		"title":   post.Title,
	})

	utils.SuccessResponse(c, http.StatusOK, post)
}

// DeletePost handles DELETE /posts/:id
func (pc *PostController) DeletePost(c *gin.Context) {
	idParam := c.Param("id")
	utils.LogInfo("Deleting post", utils.LogFields{
		"post_id": idParam,
	})

	postID, err := uuid.Parse(idParam)
	if err != nil {
		utils.LogError("Invalid post ID format for deletion", err, utils.LogFields{
			"post_id": idParam,
		})
		utils.ValidationErrorResponse(c, "Invalid post ID format")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.LogError("User not authenticated for post deletion", err, utils.LogFields{
			"post_id": postID,
		})
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	err = pc.postService.DeletePost(postID, userID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			utils.LogError("Post not found for deletion", err, utils.LogFields{
				"post_id": postID,
				"user_id": userID,
			})
			utils.NotFoundResponse(c, "Post")
			return
		}
		if utils.IsForbiddenError(err) {
			utils.LogError("User not authorized to delete post", err, utils.LogFields{
				"post_id": postID,
				"user_id": userID,
			})
			utils.ForbiddenResponse(c, "You can only delete your own posts")
			return
		}
		utils.LogError("Failed to delete post", err, utils.LogFields{
			"post_id": postID,
			"user_id": userID,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("Post deleted successfully", utils.LogFields{
		"post_id": postID,
		"user_id": userID,
	})

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// ListPosts handles GET /posts
func (pc *PostController) ListPosts(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	utils.LogInfo("Listing posts", utils.LogFields{
		"limit":  limitStr,
		"offset": offsetStr,
	})

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		utils.LogError("Invalid limit parameter", err, utils.LogFields{
			"limit": limitStr,
		})
		utils.ValidationErrorResponse(c, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.LogError("Invalid offset parameter", err, utils.LogFields{
			"offset": offsetStr,
		})
		utils.ValidationErrorResponse(c, "Invalid offset parameter")
		return
	}

	posts, err := pc.postService.ListPosts(limit, offset)
	if err != nil {
		utils.LogError("Failed to list posts", err, utils.LogFields{
			"limit":  limit,
			"offset": offset,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	utils.LogInfo("Posts listed successfully", utils.LogFields{
		"count":  len(posts),
		"limit":  limit,
		"offset": offset,
	})

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"posts":  posts,
		"limit":  limit,
		"offset": offset,
		"count":  len(posts),
	})
}

// ListPostsByUser handles GET /users/:userId/posts
func (pc *PostController) ListPostsByUser(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID format")
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "10")
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

	posts, err := pc.postService.ListPostsByUser(userID, limit, offset)
	if err != nil {
		utils.LogError("Failed to get posts by user", err, utils.LogFields{
			"user_id": userID,
			"limit":   limit,
			"offset":  offset,
		})
		utils.InternalServerErrorResponse(c, utils.GetErrorMessages().InternalServerError)
		return
	}

	// Convert to response format
	postResponses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = post.ToResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"posts":   postResponses,
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
		"count":   len(postResponses),
	})
}
