package handler

import (
	"github.com/KrittayotToin/quote-vote-backend/config"
	"github.com/KrittayotToin/quote-vote-backend/dto"
	"github.com/KrittayotToin/quote-vote-backend/model"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserController handles user-related operations
type UserController struct {
	DB *gorm.DB
}

// NewUserController creates a new user controller
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// Register creates a new user account
func (c *UserController) Register(ctx *fiber.Ctx) error {
	// Parse request body
	var input dto.UserStruct
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to process password",
		})
	}

	// Create user in database
	user := model.User{
		Email:        input.Email,
		FullName:     input.FullName,
		PasswordHash: string(hashedPassword),
	}

	if err := c.DB.Create(&user).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
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

	// Return success response with user info and token (exclude password_hash)
	return ctx.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "User registered successfully",
		"data": fiber.Map{
			"user": fiber.Map{
				"id":         user.ID,
				"email":      user.Email,
				"full_name":  user.FullName,
				"created_at": user.CreatedAt,
			},
			"token": token,
		},
	})
}

// GetProfile returns the current user's profile (protected route)
func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	// Get user info from JWT token (set by middleware)
	userID := ctx.Locals("user_id").(uint)

	// Find user in database
	var user model.User
	if err := c.DB.First(&user, userID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	// Return user profile
	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Profile retrieved successfully",
		"data": fiber.Map{
			"id":         user.ID,
			"email":      user.Email,
			"full_name":  user.FullName,
			"created_at": user.CreatedAt,
		},
	})
}

// GetCountUser returns the total number of users
func (c *UserController) GetCountUser(ctx *fiber.Ctx) error {
	var count int64

	// Count all users in database
	if err := c.DB.Model(&model.User{}).Count(&count).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to count users",
		})
	}

	// Return success response with count
	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "User count retrieved successfully",
		"data": fiber.Map{
			"total_users": count,
		},
	})
}
