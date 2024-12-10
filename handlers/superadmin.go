package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	logger2 "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MakeSuperadminData struct {
	UserID      uint                   `json:"user_id"`
	Permissions []model.PermissionName `json:"permissions"`
}

type LimitUserQuestionnaireCountData struct {
	UserID uint `json:"user_id"`
	Max    int  `json:"max"`
}

func MakeSuperadmin(superadminService service.SuperadminService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogSuperAdmin))
		user := c.Locals("user")
		if user == nil {
			logger.Error("unauthorized")
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
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusUnprocessableEntity, "failed to process entity", nil)
		}

		err = superadminService.MakeSuperadmin(loggedInUser.ID, data.UserID, data.Permissions)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to make super admin", err.Error())
		}

		logger.Info("super admin created")
		return response.Success(c, fiber.StatusCreated, "superadmin created", nil)
	}
}

func LimitUserQuestionnaireCount(superadminService service.SuperadminService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogSuperAdmin))
		data := new(LimitUserQuestionnaireCountData)
		err := c.BodyParser(data)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusUnprocessableEntity, "data is no valid", nil)
		}

		if data.Max == 0 {
			logger.Error("max value is not valid")
			return response.Error(c, fiber.StatusUnprocessableEntity, "max value is not valid", nil)
		}

		err = superadminService.LimitUserQuestionnairesCount(data.UserID, data.Max)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong with saving limit user questionnaires count", nil)
		}

		logger.Info("limit questionnaire count successfully")
		return response.Success(c, fiber.StatusCreated, "questionnaires count limitation saved", nil)
	}
}
