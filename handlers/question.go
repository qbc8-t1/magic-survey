package handlers

import (
	"strconv"

	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	logger2 "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"go.uber.org/zap"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func CreateQuestionHandler(svc service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestion))

		var questionDTO model.CreateQuestionDTO
		if err := c.BodyParser(&questionDTO); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err := questionDTO.Validate()

		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		err = svc.CreateQuestion(&questionDTO)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("question created")
		return response.Success(c, fiber.StatusCreated, "question created successfully", nil)
	}
}

func GetQuestionHandler(svc service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestion))

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		res, err := svc.GetQuestionByID(model.QuestionID(id))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("question Found")
		return response.Success(c, fiber.StatusOK, "question found successfully", res)
	}
}

func GetQuestionsByQuestionnaireIDHandler(svc service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestion))

		questionnaireIdStr := c.Params("questionnaire_id")
		questionnaireId, err := strconv.Atoi(questionnaireIdStr)
		if err != nil || questionnaireId <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		questions, err := svc.GetQuestionsByQuestionnaireID(model.QuestionnaireID(questionnaireId))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("questions retrieved")
		return response.Success(c, fiber.StatusOK, "questions retrieved successfully", questions)
	}
}

func UpdateQuestionHandler(svc service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestion))

		owner, ok := c.Locals("user").(model.User)
		if !ok {
			logger.Error("unauthorized")
			return response.Error(c, fiber.StatusBadRequest, "unauthorized", nil)
		}

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a positive integer", nil)
		}

		var questionDTO model.UpdateQuestionDTO
		if err := c.BodyParser(&questionDTO); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err = questionDTO.Validate()
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		err = svc.UpdateQuestion(model.QuestionID(id), owner.ID, &questionDTO)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("question updated")
		return response.Success(c, fiber.StatusOK, "question updated successfully", nil)
	}
}

func DeleteQuestionHandler(svc service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestion))

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		err = svc.DeleteQuestion(model.QuestionID(id))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("question deleted")
		return response.Success(c, fiber.StatusOK, "question deleted successfully", nil)
	}
}
