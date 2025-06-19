package main

import (
	"os"

	"github.com/KrittayotToin/quote-vote-backend/config"
	"github.com/KrittayotToin/quote-vote-backend/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connect to database
	db := config.ConnectDatabase()

	// Create Fiber app
	app := fiber.New()

	// Enable CORS
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Create handlers (controllers)
	userController := handler.NewUserController(db)
	authController := handler.NewAuthController(db)

	quoteController := handler.NewQuoteController(db)

	// Public routes (no authentication required)
	app.Post("/api/v1/register", userController.Register)
	app.Post("/api/v1/login", authController.Login)

	// Protected routes (authentication required)
	protected := app.Group("/api/v1")
	protected.Use(handler.AuthMiddleware())
	protected.Get("/profile", userController.GetProfile)
	protected.Get("/users/count", userController.GetCountUser)

	// Quote
	protected.Post("/quotes", quoteController.Create)
	protected.Get("/quotes", quoteController.GetAllQuotes)
	protected.Post("/quotes/:id/vote", quoteController.VoteQuote)
	protected.Get("/quotes/:id/votes", quoteController.GetVotesForQuote)
	protected.Put("/quotes/:id", quoteController.UpdateQuote)

	// Start server
	app.Listen(":8080")
}
