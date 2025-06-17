package handler

import (
	iface "github.com/KrittayotToin/quote-vote-backend/interfaces"
	"github.com/KrittayotToin/quote-vote-backend/model"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService iface.UserInterface
}

func NewUserHandler(userService iface.UserInterface) *UserHandler {
	return &UserHandler{UserService: userService}
}

// func checkPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Cannot parse JSON",
			"data":    nil,
		})
	}

	createdUser, err := h.UserService.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to create user",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "User created successfully",
		"data":    createdUser,
	})
}
