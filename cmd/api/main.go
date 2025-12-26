package main

import (
	"construction-backend/internal/config"
	"construction-backend/internal/database"
	"construction-backend/internal/handlers"
	"construction-backend/internal/middleware"
	"construction-backend/internal/repository"
	"construction-backend/internal/service"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	workerRepo := repository.NewWorkerRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, workerRepo, customerRepo, cfg.JWTSecret, db.DB)
	workerService := service.NewWorkerService(workerRepo, userRepo)
	customerService := service.NewCustomerService(customerRepo)
	bookingService := service.NewBookingService(bookingRepo, customerRepo)
	reviewService := service.NewReviewService(reviewRepo)
	messageService := service.NewMessageService(messageRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	workerHandler := handlers.NewWorkerHandler(workerService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	messageHandler := handlers.NewMessageHandler(messageService)

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	allowedOrigins := cfg.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Construction API is running"})
	})

	// Serve uploaded files
	router.Static("/uploads", "./uploads")

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.GetCurrentUser)
			auth.PUT("/profile/:userId", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.UpdateProfile)
				auth.POST("/profile/:userId/avatar", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.UploadAvatar)
			}

			// Worker routes
			workers := api.Group("/workers")
			{
				workers.GET("", workerHandler.GetAllWorkers)
				workers.GET("/user/:userId", middleware.AuthMiddleware(cfg.JWTSecret), workerHandler.GetWorkerByUserID)
			workers.GET("/filter", workerHandler.FilterWorkers)
		// Get worker by ID
		workers.GET("/:id", workerHandler.GetWorkerByID)
		workers.PUT("/:id", middleware.AuthMiddleware(cfg.JWTSecret), workerHandler.UpdateWorker)
		workers.POST("/:id/portfolio", middleware.AuthMiddleware(cfg.JWTSecret), workerHandler.AddPortfolio)
		workers.GET("/:id/reviews", reviewHandler.GetWorkerReviews)
		workers.POST("/:id/reviews", middleware.AuthMiddleware(cfg.JWTSecret), reviewHandler.CreateReview)
		}

		// Customer routes
		customers := api.Group("/customers")
		{
			customers.GET("/:id", middleware.AuthMiddleware(cfg.JWTSecret), customerHandler.GetCustomerProfile)
			customers.PUT("/:id", middleware.AuthMiddleware(cfg.JWTSecret), customerHandler.UpdateCustomerProfile)
			customers.GET("/:id/bookings", middleware.AuthMiddleware(cfg.JWTSecret), customerHandler.GetBookingHistory)
			customers.GET("/:id/favorites", middleware.AuthMiddleware(cfg.JWTSecret), customerHandler.GetFavoriteWorkers)
			customers.POST("/:id/favorites", middleware.AuthMiddleware(cfg.JWTSecret), customerHandler.AddFavoriteWorker)
			customers.DELETE("/:id/favorites/:workerId", middleware.AuthMiddleware(cfg.JWTSecret), customerHandler.RemoveFavoriteWorker)
		}

		// Booking routes
		bookings := api.Group("/bookings")
		{
			bookings.POST("", middleware.AuthMiddleware(cfg.JWTSecret), bookingHandler.CreateBooking)
			bookings.GET("/user/:userId", middleware.AuthMiddleware(cfg.JWTSecret), bookingHandler.GetUserBookings)
			bookings.GET("/worker/:workerId", middleware.AuthMiddleware(cfg.JWTSecret), bookingHandler.GetWorkerBookings)
			bookings.PATCH("/:id/status", middleware.AuthMiddleware(cfg.JWTSecret), bookingHandler.UpdateBookingStatus)
			bookings.DELETE("/:id", middleware.AuthMiddleware(cfg.JWTSecret), bookingHandler.CancelBooking)
		}

		// Message routes
		messages := api.Group("/messages")
		{
			messages.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			messages.POST("", messageHandler.SendMessage)
			messages.GET("/conversations", messageHandler.GetConversations)
			messages.GET("/:userId", messageHandler.GetMessages)
			messages.PATCH("/:id/read", messageHandler.MarkAsRead)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}