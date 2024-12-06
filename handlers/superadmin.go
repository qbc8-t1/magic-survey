package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type MakeSuperadminData struct {
	UserID      uint     `json:"user_id"`
	Permissions []string `json:"permissions"`
}

type LimitUserQuestionnaireCountData struct {
	UserID uint `json:"user_id"`
	Max    int  `json:"max"`
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

func LimitUserQuestionnaireCount(superadminService service.SuperadminService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := new(LimitUserQuestionnaireCountData)
		err := c.BodyParser(data)
		if err != nil {
			return response.Error(c, fiber.StatusUnprocessableEntity, "data is no valid", err.Error())
		}

		if data.Max == 0 {
			return response.Error(c, fiber.StatusUnprocessableEntity, "max value is not valid", nil)
		}

		err = superadminService.LimitUserQuestionnairesCount(data.UserID, data.Max)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong with saving limit user questionnaires count", err.Error())
		}

		return c.Status(fiber.StatusCreated).SendString("limit user questionnaire count saved successfully")
	}
}
