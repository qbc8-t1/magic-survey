package handlers

import (
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	logger2 "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func QuestionnaireCreate(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestionnaire))

		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("user variable was't in locals")
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			logger.Error("failed to parse locals user")
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		can, err := qService.CheckIfUserCanMakeNewQuestionnaire(user)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusForbidden, "something went wrong to check the questionnaires count", err)
		}
		if !can {
			logger.Error("user reached the limitation to make questionnaire")
			return response.Error(c, fiber.StatusForbidden, "you have reached your limitation to make questionnaires", nil)
		}

		var createData model.CreateQuestionnaireDTO
		if err := c.BodyParser(&createData); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		questionnaireRawObject, err := createData.ValidateAndMakeObjectForCreate()
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		questionnaireRawObject.OwnerID = model.UserID(user.ID)
		questionnaire, err := qService.CreateQuestionnaire(&questionnaireRawObject)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong with creating new questionnaire", err)
		}

		logger.Info("questinonaire created")
		return response.Success(c, fiber.StatusCreated, "Questionnaire Created Successfully", model.ToQuestionnaireResponse(&questionnaire))
	}
}

func GetQuestionnairesList(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestionnaire))

		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("failed to get user from locals")
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			logger.Error("failed to cast locals user")
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		page := c.QueryInt("page")
		if page == 0 {
			page = 1
		}

		qList, err := qService.GetQuestionnairesList(model.UserID(user.ID), page)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get questionnaires list", nil)
		}

		if len(qList) == 0 {
			if page == 1 {
				logger.Info("no questionnaires found")
				return response.Success(c, fiber.StatusOK, "you don't have any questionnaires yet", nil)
			} else {
				logger.Info("no questionnaires found for this page")
				return response.Success(c, fiber.StatusOK, "no items for this page", nil)
			}
		}

		logger.Info("list of questionnaires")
		return response.Success(c, fiber.StatusOK, "list of your questionnaires", qList)
	}
}

func QuestionnaireUpdate(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestionnaire))

		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("failed to get user from locals")
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		_, ok := localUser.(model.User)
		if !ok {
			logger.Error("failed to cast locals user")
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id param is invalid", nil)
		}

		questionnaire, err := qService.GetQuestionnaireByID(model.QuestionnaireID(questionnaireID))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "failed to get the questionnaire", nil)
		}

		var updateData model.UpdateQuestionnaireDTO
		if err := c.BodyParser(&updateData); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		// we need these data to check

		if updateData.CanSubmitFrom == "" {
			updateData.CanSubmitFrom = questionnaire.CanSubmitFrom.Format(time.DateTime)
		}

		if updateData.CanSubmitUntil == "" {
			updateData.CanSubmitUntil = questionnaire.CanSubmitUntil.Format(time.DateTime)
		}

		if updateData.MaxMinutesToResponse == nil {
			updateData.MaxMinutesToResponse = &questionnaire.MaxMinutesToResponse
		}

		questionnaireRawObject, err := updateData.ValidateAndMakeObjectForUpdate()
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		err = qService.UpdateQuestionaire(model.QuestionnaireID(questionnaireID), &questionnaireRawObject)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to update questionnaire", nil)
		}

		logger.Info("questionnaire updated")
		return response.Success(c, fiber.StatusCreated, "questionnaire updated", nil)
	}
}

func QuestionnaireGet(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestionnaire))

		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("failed to get user from locals")
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		_, ok := localUser.(model.User)
		if !ok {
			logger.Error("failed to cast locals user")
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id param is invalid", nil)
		}

		questionnaire, err := qService.GetQuestionnaireByID(model.QuestionnaireID(questionnaireID))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get the questionnaire", nil)
		}

		logger.Info("questionnaire found")
		return response.Success(c, fiber.StatusOK, "questionnaire data", model.ToQuestionnaireResponse(questionnaire))
	}
}

func QuestionnaireDelete(qService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestionnaire))

		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("failed to get user from locals")
			return response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		_, ok := localUser.(model.User)
		if !ok {
			logger.Error("failed to cast locals user")
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id param is invalid", nil)
		}

		err = qService.DeleteQuestionnaire(model.QuestionnaireID(questionnaireID))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to delete the questionnaire", nil)
		}

		logger.Info("questionnaire delete")
		return response.Success(c, fiber.StatusOK, "questionnaire deleted", nil)
	}
}
