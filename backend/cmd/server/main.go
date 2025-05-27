package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/readitlater/backend/internal/config"
	"github.com/readitlater/backend/internal/database"
	"github.com/readitlater/backend/internal/auth"
	"github.com/readitlater/backend/internal/handlers"
	"github.com/readitlater/backend/internal/middleware"
	"github.com/readitlater/backend/internal/services"
	"github.com/readitlater/backend/internal/storage"
	"github.com/readitlater/backend/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg)
	logger.Info("Starting Read-It-Later backend server")

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize Redis client
	redisClient, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize services
	authService := auth.NewService(cfg)
	storageService := storage.NewService(cfg, logger)
	
	userService := services.NewUserService(db, authService, logger)
	articleService := services.NewArticleService(db, storageService, logger)
	captureService := services.NewCaptureService(cfg, articleService, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, authService, logger)
	userHandler := handlers.NewUserHandler(userService, logger)
	articleHandler := handlers.NewArticleHandler(articleService, logger)
	captureHandler := handlers.NewCaptureHandler(captureService, logger)

	// Setup Gin router
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.RateLimit(cfg, redisClient))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(authService, userService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", userHandler.GetProfile)
				users.PUT("/me", userHandler.UpdateProfile)
				users.POST("/change-password", userHandler.ChangePassword)
				users.GET("/stats", userHandler.GetStats)
			}

			// Article routes
			articles := protected.Group("/articles")
			{
				articles.GET("", articleHandler.GetArticles)
				articles.POST("", articleHandler.CreateArticle)
				articles.GET("/:id", articleHandler.GetArticle)
				articles.PUT("/:id", articleHandler.UpdateArticle)
				articles.DELETE("/:id", articleHandler.DeleteArticle)
				articles.GET("/:id/content", articleHandler.GetArticleContent)
				articles.POST("/:id/read", articleHandler.MarkAsRead)
				articles.POST("/:id/favorite", articleHandler.ToggleFavorite)
				articles.POST("/:id/archive", articleHandler.ToggleArchive)
			}

			// Capture routes
			capture := protected.Group("/capture")
			{
				capture.POST("/url", captureHandler.CaptureURL)
				capture.GET("/status/:id", captureHandler.GetCaptureStatus)
			}

			// Search routes
			search := protected.Group("/search")
			{
				search.GET("/articles", articleHandler.SearchArticles)
				search.GET("/suggestions", articleHandler.GetSearchSuggestions)
			}

			// Export/Import routes
			export := protected.Group("/export")
			{
				export.GET("/articles", articleHandler.ExportArticles)
				export.POST("/import", articleHandler.ImportArticles)
			}
		}

		// Extension API routes (with API key authentication)
		extension := api.Group("/extension")
		extension.Use(middleware.ExtensionAuth(cfg))
		{
			extension.POST("/capture", captureHandler.ExtensionCapture)
			extension.GET("/user", userHandler.GetExtensionUser)
		}
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
} 