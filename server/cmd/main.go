package main

import (
	"os"

	"github.com/TejasThombare20/post-comments-service/config"
	"github.com/TejasThombare20/post-comments-service/controllers"
	"github.com/TejasThombare20/post-comments-service/middleware"
	"github.com/TejasThombare20/post-comments-service/repository"
	"github.com/TejasThombare20/post-comments-service/routes"
	"github.com/TejasThombare20/post-comments-service/services"
	"github.com/TejasThombare20/post-comments-service/utils"
	"github.com/TejasThombare20/post-comments-service/validator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger first
	utils.InitLogger()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		utils.LogWarn("No .env file found, using system environment variables", nil)
	}

	// Initialize database connection
	db, err := config.InitDB()
	if err != nil {
		utils.LogError("Failed to connect to database", err, nil)
		os.Exit(1)
	}

	// Initialize validator
	validator := validator.NewValidator()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo, userRepo)
	commentService := services.NewCommentService(commentRepo, postRepo, userRepo, validator)
	jwtService := services.NewJWTService()
	authService := services.NewAuthService(userRepo, jwtService, userService, validator)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	postController := controllers.NewPostController(postService)
	commentController := controllers.NewCommentController(commentService)
	authController := controllers.NewAuthController(authService, validator)

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(router, userController, postController, commentController, authController, jwtService)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.LogInfo("Server starting", utils.LogFields{"port": port})
	if err := router.Run(":" + port); err != nil {
		utils.LogError("Failed to start server", err, utils.LogFields{"port": port})
		os.Exit(1)
	}
}
