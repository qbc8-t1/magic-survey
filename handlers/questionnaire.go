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
		// TODO - get logged in user
		// TODO - check if user can make a new questionnaire - check if it has reached the limitation or not

		var requestData model.CreateQuestionnaireDTO
		if err := c.BodyParser(&requestData); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		questionnaireRawObject, err := requestData.ValidateAndMakeObject()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		// TODO - put user id in the questionnaireRawObject
		questionnaireRawObject.OwnerID = 2

		questionnaire, err := qService.CreateQuestionnaire(&questionnaireRawObject)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong with creating new questionnaire", err)
		}

		return response.Success(c, fiber.StatusCreated, "Questionnaire Created Successfully", model.ToQuestionnaireResponse(&questionnaire))
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
