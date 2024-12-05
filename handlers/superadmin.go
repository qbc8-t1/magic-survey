package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func LimitUserQuestionnaireCount(superadminService *service.SuperadminService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

type MakeSuperadminData struct {
	UserID      uint     `json:"user_id"`
	Permissions []string `json:"permissions"`
}

func MakeSuperadmin(superadminService service.SuperadminService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			// If the user is not set, return unauthorized
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Type assert the user to your User struct
		loggedInUser, ok := user.(model.User)
		if !ok {
			// Handle the case where the type assertion fails
			return response.Error(c, fiber.StatusInternalServerError, "failed to get auth user", nil)
		}

		data := new(MakeSuperadminData)
		err := c.BodyParser(&data)
		if err != nil {
			return response.Error(c, fiber.StatusUnprocessableEntity, "failed to process entity", err.Error())
		}

		err = superadminService.MakeSuperadmin(loggedInUser.ID, data.UserID, data.Permissions)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to make superadmin", err.Error())
		}

		return c.Status(fiber.StatusCreated).SendString("superadmin created")
	}
}
