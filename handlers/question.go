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

func CreateQuestionHandler(service service.IQuestionService) fiber.Handler {
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

		err = service.CreateQuestion(&questionDTO)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("question created")
		return response.Success(c, fiber.StatusCreated, "question created successfully", nil)
	}
}

func GetQuestionHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogQuestion))

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		res, err := service.GetQuestionByID(model.QuestionID(id))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("question Found")
		return response.Success(c, fiber.StatusOK, "question found successfully", res)
	}
}

func GetQuestionsByQuestionnaireIDHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		questionnaireIdStr := c.Params("questionnaire_id")
		questionnaireId, err := strconv.Atoi(questionnaireIdStr)
		if err != nil || questionnaireId <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		questions, err := service.GetQuestionsByQuestionnaireID(model.QuestionnaireID(questionnaireId))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "questions retrieved successfully", questions)
	}
}

func UpdateQuestionHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a positive integer", nil)
		}

		var questionDTO model.UpdateQuestionDTO
		if err := c.BodyParser(&questionDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err = questionDTO.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		err = service.UpdateQuestion(model.QuestionID(id), &questionDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "question updated successfully", nil)
	}
}

func DeleteQuestionHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		err = service.DeleteQuestion(model.QuestionID(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "question deleted successfully", nil)
	}
}
