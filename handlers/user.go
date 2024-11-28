package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/gofiber/fiber/v2"
)

// 1. +get data
// 2. +init db connection
// 3. password hashing (add salt row in db)
// 4. signup
// 5. login
// 6. setup SMTP server for sending email
// 7. send verification code

func UserCreate(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dto model.CreateUserDTO
		if err := c.BodyParser(&dto); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		user := model.ToUserModel(&dto)

		createdUser, err := userService.CreateUser(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(model.ToUserResponse(createdUser))
	}
}
