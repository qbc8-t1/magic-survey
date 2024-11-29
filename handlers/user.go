package handlers

import (
	"errors"
	"fmt"
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/jwt"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// 5. login
// 6. setup SMTP server for sending email
// 7. send verification code
// 8. handle 2 step verification

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
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "User Created", model.ToUserResponse(createdUser))
	}
}

func Login(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", nil)
		}

		user, err := userService.LoginUser(&req)
		fmt.Println(err)

		if errors.Is(err, service.ErrWrongEmailPass) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		token, err := jwt.GenerateToken(uint(user.ID))

		c.Cookie(&fiber.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  jwt.GetTokenExpiry(), // Token expiry time
			Secure:   true,
			HTTPOnly: true,
			SameSite: fiber.CookieSameSiteStrictMode,
		})

		return response.Success(c, fiber.StatusOK, "Login successful", model.ToUserResponse(user))
	}
}
