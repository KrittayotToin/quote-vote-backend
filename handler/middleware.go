package handler

import (
	"strings"

	"github.com/KrittayotToin/quote-vote-backend/config"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware protects routes that require authentication
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Authorization header required",
			})
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization format",
			})
		}

		// Extract token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := config.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		// Store user info in context for later use
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		// Continue to next handler
		return c.Next()
	}
}
