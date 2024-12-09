package handlers

import (
	"errors"
	"strconv"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	logger2 "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// StartHandler starts a questionnaire
func StartHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))

		// Extract user from context
		user, ok := c.Locals("user").(model.User)
		if !ok || user.ID <= 0 {
			logger.Error("userID is required and must be greater than 0")
			return response.Error(c, fiber.StatusBadRequest, "userID is required and must be greater than 0", nil)
		}

		questionnaireIDStr := c.Params("questionnaire_id")
		questionnaireID, err := strconv.ParseUint(questionnaireIDStr, 10, 64)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid questionnaire_id", nil)
		}

		question, err := svc.Start(model.QuestionnaireID(questionnaireID), model.UserID(user.ID))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "questionnaire started", question)
	}
}

func SubmitHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))

		// Extract user from context
		user, ok := c.Locals("user").(model.User)
		if !ok || user.ID <= 0 {
			logger.Error("userID is required and must be greater than 0")
			return response.Error(c, fiber.StatusBadRequest, "userID is required and must be greater than 0", nil)
		}

		// Parse questionID from params
		questionIDStr := c.Params("question_id")
		questionID, err := strconv.ParseUint(questionIDStr, 10, 64)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid question_id", nil)
		}

		// Parse the answer DTO from the request body
		var dto model.CreateAnswerDTO
		if err := c.BodyParser(&dto); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
		}

		// Validate the DTO
		err = dto.Validate()
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		// Map DTO to Answer model
		answer := model.ToAnswerModel(&dto)

		// Call the service to submit the answer
		err = svc.Submit(model.QuestionID(questionID), answer, model.UserID(user.ID))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		logger.Info("answer submited")
		return response.Success(c, fiber.StatusOK, "answer submitted successfully", nil)
	}
}

// BackHandler moves user to the previous question
func BackHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))

		// Extract user from context
		user, ok := c.Locals("user").(model.User)
		if !ok || user.ID <= 0 {
			logger.Error("userID is required and must be greater than 0")
			return response.Error(c, fiber.StatusBadRequest, "userID is required and must be greater than 0", nil)
		}

		questionResponse, err := svc.Back(model.UserID(user.ID))
		if err != nil {
			if errors.Is(err, service.ErrCannotGoBack) {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
			}

			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		logger.Info("previous question")
		return response.Success(c, fiber.StatusOK, "previous question", questionResponse)
	}
}

// NextHandler moves user to the next question
func NextHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))

		// Extract user from context
		user, ok := c.Locals("user").(model.User)
		if !ok || user.ID <= 0 {
			logger.Error("userID is required and must be greater than 0")
			return response.Error(c, fiber.StatusBadRequest, "userID is required and must be greater than 0", nil)
		}

		questionResponse, err := svc.Next(model.UserID(user.ID))
		if err != nil {
			if errors.Is(err, service.ErrNoNextQuestion) {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusNotFound, "no more questions left", nil)
			}

			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		logger.Info("next question")
		return response.Success(c, fiber.StatusOK, "next question", questionResponse)
	}
}

// EndHandler finalizes the questionnaire submission
func EndHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))

		// Extract user from context
		user, ok := c.Locals("user").(model.User)
		if !ok || user.ID <= 0 {
			logger.Error("userID is required and must be greater than 0")
			return response.Error(c, fiber.StatusBadRequest, "userID is required and must be greater than 0", nil)
		}

		err := svc.End(model.UserID(user.ID))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		logger.Info("questionnaire ended")
		return response.Success(c, fiber.StatusOK, "questionnaire ended", nil)
	}
}
