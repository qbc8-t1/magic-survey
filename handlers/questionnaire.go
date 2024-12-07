package handlers

import (
	"strconv"
	"time"

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

		var createData model.CreateQuestionnaireDTO
		if err := c.BodyParser(&createData); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		questionnaireRawObject, err := createData.ValidateAndMakeObjectForCreate()
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

func GetQuestionnairesList(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		page := c.QueryInt("page")
		if page == 0 {
			page = 1
		}

		qList, err := qService.GetQuestionnairesList(model.UserID(user.ID), page)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get questionnaires list", err.Error())
		}

		if len(qList) == 0 {
			return response.Success(c, fiber.StatusOK, "you don't have any questionnaires yet", nil)
		}

		return response.Success(c, fiber.StatusOK, "list of your questionnaires", qList)
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

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id param is invalid", nil)
		}

		questionnaire, err := qService.GetQuestionnaireByID(model.QuestionnaireID(questionnaireID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "failed to get the questionnaire", err.Error())
		}

		var updateData model.UpdateQuestionnaireDTO
		if err := c.BodyParser(&updateData); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		if updateData.CanSubmitFrom == "" {
			updateData.CanSubmitFrom = questionnaire.CanSubmitFrom.Format(time.DateTime)
		}

		if updateData.CanSubmitUntil == "" {
			updateData.CanSubmitUntil = questionnaire.CanSubmitUntil.Format(time.DateTime)
		}

		if updateData.MaxMinutesToResponse == "" {
			updateData.MaxMinutesToResponse = strconv.Itoa(questionnaire.MaxMinutesToResponse)
		}

		questionnaireRawObject, err := updateData.ValidateAndMakeObjectForUpdate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		err = qService.UpdateQuestionaire(model.QuestionnaireID(questionnaireID), &questionnaireRawObject)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to update questionnaire", err.Error())
		}

		return response.Success(c, fiber.StatusCreated, "questionnaire updated", nil)
	}
}

func QuestionnaireGet(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		_, ok := localUser.(model.User)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id param is invalid", nil)
		}

		questionnaire, err := qService.GetQuestionnaireByID(model.QuestionnaireID(questionnaireID))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get the questionnaire", err.Error())
		}

		return response.Success(c, fiber.StatusOK, "questionnaire data", model.ToQuestionnaireResponse(&questionnaire))
	}
}

func QuestionnaireDelete(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		_, ok := localUser.(model.User)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id param is invalid", nil)
		}

		err = qService.DeleteQuestionnaire(model.QuestionnaireID(questionnaireID))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to delete the questionnaire", err.Error())
		}

		return response.Success(c, fiber.StatusOK, "questionnaire deleted", nil)
	}
}
