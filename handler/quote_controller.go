package handler

import (
	"log"

	"github.com/KrittayotToin/quote-vote-backend/dto"
	"github.com/KrittayotToin/quote-vote-backend/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuoteController struct {
	DB *gorm.DB
}

func NewQuoteController(db *gorm.DB) *QuoteController {
	return &QuoteController{DB: db}
}

func (c *QuoteController) Create(ctx *fiber.Ctx) error {
	var body dto.QuoteStruct

	// Log the incoming request
	log.Printf("Creating quote - Request received")

	if err := ctx.BodyParser(&body); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
		})
	}

	// Get user ID from JWT token
	userID := ctx.Locals("user_id").(uint)

	// Convert DTO to model
	quote := model.Quote{
		Text:      body.Text,
		Author:    body.Author,
		Votes:     0, // Default votes
		CreatedBy: userID,
	}

	// Create quote in database
	if err := c.DB.Create(&quote).Error; err != nil {
		log.Printf("Error creating quote: %v", err)
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create Quote",
		})
	}

	// Log the created quote
	log.Printf("Quote created successfully: %+v", quote)

	return ctx.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Quote created successfully",
		"data":    quote,
	})
}

func (c *QuoteController) GetAllQuotes(ctx *fiber.Ctx) error {
	var quotes []model.Quote

	// Log the request
	log.Printf("Getting all quotes - Request received")

	// Find all quotes in database
	if err := c.DB.Find(&quotes).Error; err != nil {
		log.Printf("Error retrieving quotes: %v", err)
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve quotes",
		})
	}

	// Log the result
	log.Printf("Retrieved %d quotes successfully", len(quotes))

	// Format quotes to match TypeScript interface
	var formattedQuotes []fiber.Map
	for _, quote := range quotes {
		formattedQuotes = append(formattedQuotes, fiber.Map{
			"id":         quote.ID,
			"text":       quote.Text,
			"author":     quote.Author,
			"votes":      quote.Votes,
			"created_at": quote.CreatedAt,
			"created_by": quote.CreatedBy,
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Quotes retrieved successfully",
		"data": fiber.Map{
			"quotes": formattedQuotes,
			"count":  len(quotes),
		},
	})
}

func (c *QuoteController) VoteQuote(ctx *fiber.Ctx) error {
	// Get quote ID from URL parameter
	quoteID := ctx.Params("id")
	userID := ctx.Locals("user_id").(uint)

	// Log the vote request
	log.Printf("Voting on quote ID: %s by user ID: %d", quoteID, userID)

	// Check if user already voted on this quote
	var existingVote model.Vote
	if err := c.DB.Where("user_id = ? AND quote_id = ?", userID, quoteID).First(&existingVote).Error; err == nil {
		// User already voted
		log.Printf("User %d already voted on quote %s", userID, quoteID)
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "You have already voted on this quote",
		})
	}

	// Get the quote to update votes
	var quote model.Quote
	if err := c.DB.First(&quote, quoteID).Error; err != nil {
		log.Printf("Quote not found: %s", quoteID)
		return ctx.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Quote not found",
		})
	}

	// Start transaction
	tx := c.DB.Begin()

	// Create vote record
	vote := model.Vote{
		UserID:  userID,
		QuoteID: quote.ID,
	}

	if err := tx.Create(&vote).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating vote: %v", err)
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create vote",
		})
	}

	// Update quote votes count
	if err := tx.Model(&quote).Update("votes", quote.Votes+1).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating quote votes: %v", err)
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update quote votes",
		})
	}

	// Commit transaction
	tx.Commit()

	log.Printf("Vote successful - Quote %s now has %d votes", quoteID, quote.Votes+1)

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Vote recorded successfully",
		"data": fiber.Map{
			"quote_id": quote.ID,
			"votes":    quote.Votes + 1,
		},
	})
}

func (c *QuoteController) GetVotesForQuote(ctx *fiber.Ctx) error {
	quoteID := ctx.Params("id")
	var votes []model.Vote

	// Find all votes for the given quote
	if err := c.DB.Where("quote_id = ?", quoteID).Find(&votes).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve votes",
		})
	}

	var formattedVotes []fiber.Map
	for _, v := range votes {
		formattedVotes = append(formattedVotes, fiber.Map{
			"id":         v.ID,
			"quote_id":   v.QuoteID,
			"user_id":    v.UserID,
			"created_at": v.CreatedAt.Format("02 Jan 2006 15:04"),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    formattedVotes,
	})
}

func (c *QuoteController) UpdateQuote(ctx *fiber.Ctx) error {
	quoteID := ctx.Params("id")
	var body dto.QuoteStruct

	// Parse request body
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request data",
		})
	}

	// Find the quote
	var quote model.Quote
	if err := c.DB.First(&quote, quoteID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Quote not found",
		})
	}

	// Update fields
	quote.Text = body.Text
	quote.Author = body.Author

	if err := c.DB.Save(&quote).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update quote",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Quote updated successfully",
		"data": fiber.Map{
			"id":         quote.ID,
			"text":       quote.Text,
			"author":     quote.Author,
			"votes":      quote.Votes,
			"created_at": quote.CreatedAt,
		},
	})
}
