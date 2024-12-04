package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func StartHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func SubmitHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func BackHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from context
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		// Call the service method
		questionResponse, err := service.Back(model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "previous question", questionResponse)
	}
}

func NextHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from context
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		// Call the service method
		questionResponse, err := service.Next(model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "next question", questionResponse)
	}
}

func EndHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
