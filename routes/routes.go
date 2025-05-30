package routes

import (
	"github.com/TejasThombare20/post-comments-service/controllers"
	"github.com/TejasThombare20/post-comments-service/middleware"
	"github.com/TejasThombare20/post-comments-service/services"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(
	router *gin.Engine,
	userController *controllers.UserController,
	postController *controllers.PostController,
	commentController *controllers.CommentController,
	authController *controllers.AuthController,
	jwtService *services.JWTService,
) {
	// Add CORS middleware
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Post-Comments Service is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)                                                     // POST /api/v1/auth/register
			auth.POST("/login", authController.Login)                                                           // POST /api/v1/auth/login
			auth.POST("/refresh", authController.RefreshToken)                                                  // POST /api/v1/auth/refresh
			auth.POST("/logout", authController.Logout)                                                         // POST /api/v1/auth/logout
			auth.GET("/profile", middleware.AuthMiddleware(jwtService), authController.GetProfile)              // GET /api/v1/auth/profile
			auth.POST("/change-password", middleware.AuthMiddleware(jwtService), authController.ChangePassword) // POST /api/v1/auth/change-password
		}

		// Public user routes (read-only)
		users := v1.Group("/users")
		{
			users.GET("", userController.ListUsers)                            // GET /api/v1/users
			users.GET("/username/:username", userController.GetUserByUsername) // GET /api/v1/users/username/:username
			users.GET("/:userId/posts", postController.ListPostsByUser)        // GET /api/v1/users/:userId/posts
			users.GET("/user/:id", userController.GetUser)                     // GET /api/v1/users/:id
		}

		// Protected user routes (require authentication)
		protectedUsers := v1.Group("/users")
		protectedUsers.Use(middleware.AuthMiddleware(jwtService))
		{
			protectedUsers.PUT("/user/:id", userController.UpdateUser)    // PUT /api/v1/users/:id
			protectedUsers.DELETE("/user/:id", userController.DeleteUser) // DELETE /api/v1/users/:id
		}

		// Public post routes (read-only)
		posts := v1.Group("/posts")
		{
			posts.GET("", postController.ListPosts)                                   // GET /api/v1/posts
			posts.GET("/post/:id", postController.GetPost)                            // GET /api/v1/posts/:id
			posts.GET("/post/:id/comments", postController.GetPostWithComments)       // GET /api/v1/posts/:id/comments
			posts.GET("/post/:postId/comments", commentController.ListCommentsByPost) // GET /api/v1/posts/:postId/comments
		}

		// Protected post routes (require authentication)
		protectedPosts := v1.Group("/posts")
		protectedPosts.Use(middleware.AuthMiddleware(jwtService))
		{
			protectedPosts.POST("", postController.CreatePost)                             // POST /api/v1/posts
			protectedPosts.PUT("/post/:id", postController.UpdatePost)                     // PUT /api/v1/posts/:id
			protectedPosts.DELETE("/post/:id", postController.DeletePost)                  // DELETE /api/v1/posts/:id
			protectedPosts.POST("/post/:postId/comments", commentController.CreateComment) // POST /api/v1/posts/:postId/comments
		}

		// Public comment routes (read-only)
		comments := v1.Group("/comments")
		{
			comments.GET("/:id", commentController.GetComment)                // GET /api/v1/comments/:id
			comments.GET("/:id/replies", commentController.GetCommentReplies) // GET /api/v1/comments/:id/replies
		}

		// Protected comment routes (require authentication)
		protectedComments := v1.Group("/comments")
		protectedComments.Use(middleware.AuthMiddleware(jwtService))
		{
			protectedComments.PUT("/:id", commentController.UpdateComment)    // PUT /api/v1/comments/:id
			protectedComments.DELETE("/:id", commentController.DeleteComment) // DELETE /api/v1/comments/:id
		}
	}
}
