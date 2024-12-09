package handlers

import (
	"strconv"

	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	logger2 "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
)

func CreateAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))
		var answerDTO model.CreateAnswerDTO
		if err := c.BodyParser(&answerDTO); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err := answerDTO.Validate()

		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", nil)
		}

		err = service.CreateAnswer(&answerDTO)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("answer created")
		return response.Success(c, fiber.StatusCreated, "answer created successfully", nil)
	}
}

func GetAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		res, err := service.GetAnswerByID(model.AnswerID(id))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("answer Found")
		return response.Success(c, fiber.StatusOK, "answer found successfully", res)
	}
}

func UpdateAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a positive integer", nil)
		}

		var answerDTO model.UpdateAnswerDTO
		if err := c.BodyParser(&answerDTO); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err = answerDTO.Validate()
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		err = service.UpdateAnswer(model.AnswerID(id), &answerDTO)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("answer updated")
		return response.Success(c, fiber.StatusOK, "answer updated successfully", nil)
	}
}

func DeleteAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogAnswer))
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid id. the id must be a posetive integer", nil)
		}

		err = service.DeleteAnswer(model.AnswerID(id))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		logger.Info("answer deleted")
		return response.Success(c, fiber.StatusOK, "answer deleted successfully", nil)
	}
}
