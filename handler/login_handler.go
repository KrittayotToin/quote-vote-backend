package handler

import (
	"errors"

	"github.com/KrittayotToin/quote-vote-backend/config"
	"github.com/KrittayotToin/quote-vote-backend/dto"
	"github.com/KrittayotToin/quote-vote-backend/model"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthController handles authentication operations
type AuthController struct {
	DB *gorm.DB
}

// NewAuthController creates a new auth controller
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Login authenticates a user and returns a JWT token
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	// Parse request body
	var input dto.LoginStruct
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
		})
	}

	// Find user by email
	var user model.User
	if err := c.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Database error",
		})
	}

	// Check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Invalid password",
		})
	}

	// Generate JWT token
	token, err := config.GenerateToken(user.ID, user.Email)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
		})
	}

	// Return success response with token
	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}
