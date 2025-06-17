package main

import (
	"github.com/KrittayotToin/quote-vote-backend/config"
	"github.com/KrittayotToin/quote-vote-backend/handler"
	"github.com/KrittayotToin/quote-vote-backend/repository"
	"github.com/KrittayotToin/quote-vote-backend/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db := config.ConnectDatabase()
	app := fiber.New()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	app.Post("/api/v1/user", userHandler.Create)

	app.Listen(":3000")
}
