package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
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
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		user := model.ToUserModel(&dto)
		err := user.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}
		createdUser, err := userService.CreateUser(user)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "couldnt create the user", nil)
		}

		return c.Status(fiber.StatusCreated).JSON(model.ToUserResponse(createdUser))
	}
}
