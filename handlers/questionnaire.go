package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// TODO
// YOUSEF
func QuestionnaireCreate(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dto model.CreateQuestionnaireDTO
		if err := c.BodyParser(&dto); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		questionnaire := model.ToQuestionnaireModel(&dto)

		err := questionnaire.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}
		// tokens, err := userService.CreateUser(user)
		// if err != nil {
		// 	return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		// }

		// return response.Success(c, fiber.StatusCreated, "User Created", tokens)
		return nil
	}
}

func QuestionnaireUpdate(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func QuestionnaireGet(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func QuestionnaireDelete(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
