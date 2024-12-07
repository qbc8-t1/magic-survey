package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func QuestionnaireCreate(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		can, err := qService.CheckIfUserCanMakeNewQuestionnaire(user)
		if err != nil {
			return response.Error(c, fiber.StatusForbidden, "something went wrong to check the questionnaires count", err)
		}
		if !can {
			return response.Error(c, fiber.StatusForbidden, "you have reached your limitation to make questionnaires", nil)
		}

		var requestData model.CreateQuestionnaireDTO
		if err := c.BodyParser(&requestData); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		questionnaireRawObject, err := requestData.ValidateAndMakeObject()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		questionnaireRawObject.OwnerID = user.ID
		questionnaire, err := qService.CreateQuestionnaire(&questionnaireRawObject)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong with creating new questionnaire", err)
		}

		return response.Success(c, fiber.StatusCreated, "Questionnaire Created Successfully", model.ToQuestionnaireResponse(&questionnaire))
	}
}

func QuestionnaireUpdate(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		_, ok := localUser.(model.User)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		// in the rbac middleware we check ownership and permission

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "questionnaire param is invalid", nil)
		}

		var requestData model.CreateQuestionnaireDTO
		if err := c.BodyParser(&requestData); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		questionnaireRawObject, err := requestData.ValidateAndMakeObject()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		questionnaireRawObject.ID = uint(questionnaireID)

		err = qService.UpdateQuestionaire(&questionnaireRawObject)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to update questionnaire", err.Error())
		}

		return response.Success(c, fiber.StatusCreated, "questionnaire updated", nil)
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
