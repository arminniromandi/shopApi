package main

import (
	"log"
	"net/http"

	"shopProject/internal/config"
	"shopProject/internal/database"
	"shopProject/internal/handlers"
	"shopProject/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to Database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("Warning: Database connection failed (%v). Continuing to allow server setup...\n", err)
	} else {
		// Populate initial sample data (Categories & Products) if database is empty
		database.SeedData(db)
	}

	r := gin.Default()

	// CORS middleware for client/mobile requests
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Handlers initialization
	authHandler := handlers.NewAuthHandler(cfg)
	productHandler := handlers.NewProductHandler()
	cartHandler := handlers.NewCartHandler()

	api := r.Group("/api")
	{
		// Health check
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Shop API is up and running!"})
		})

		// --- Authentication Routes ---
		auth := api.Group("/auth")
		{
			auth.POST("/send-otp", authHandler.SendOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
		}

		// --- Categories & Products (Public) ---
		categories := api.Group("/categories")
		{
			categories.GET("", productHandler.GetCategories)
			categories.POST("", productHandler.CreateCategory)
			categories.GET("/:identifier/products", productHandler.GetCategoryProducts)
		}

		products := api.Group("/products")
		{
			products.GET("/special-offers", productHandler.GetSpecialOffers)
			products.GET("/:id", productHandler.GetProductDetails)
			products.POST("", productHandler.CreateProduct)
			products.GET("/:id/comments", productHandler.GetComments)

			// Protected comment creation
			products.POST("/:id/comments", middleware.AuthMiddleware(cfg), productHandler.AddComment)
		}

		// --- Cart Routes (Protected) ---
		cart := api.Group("/cart", middleware.AuthMiddleware(cfg))
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("/items", cartHandler.AddToCart)
			cart.DELETE("/items/:id", cartHandler.RemoveFromCart)
		}
	}

	log.Printf("Server starting on port :%s ...\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
