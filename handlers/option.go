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

func CreateOptionHandler(service service.IOptionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogOption))

		var optionDTO model.CreateOptionDTO
		if err := c.BodyParser(&optionDTO); err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err := optionDTO.Validate()

		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", nil)
		}

		err = service.CreateOption(&optionDTO)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "error in creating the option", nil)
		}

		logger.Info("option created")
		return response.Success(c, fiber.StatusCreated, "option created", nil)
	}
}

func GetOptionHandler(service service.IOptionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogOption))

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		res, err := service.GetOptionByID(model.OptionID(id))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "error in retrieving the option", nil)
		}

		logger.Info("option Found")
		return response.Success(c, fiber.StatusOK, "option Found", res)
	}
}

func GetOptionsByQuestionIDHandler(service service.IOptionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		questionnIdStr := c.Params("question_id")
		questionId, err := strconv.Atoi(questionnIdStr)
		if err != nil || questionId <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		options, err := service.GetOptionsByQuestionID(model.QuestionID(questionId))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to fetch options", nil)
		}

		return response.Success(c, fiber.StatusOK, "options retrieved successfully", options)
	}
}

func UpdateOptionHandler(service service.IOptionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a positive integer", nil)
		}

		var optionDTO model.UpdateOptionDTO
		if err := c.BodyParser(&optionDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err = optionDTO.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", nil)
		}

		err = service.UpdateOption(model.OptionID(id), &optionDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "error in updating the option", nil)
		}

		return response.Success(c, fiber.StatusOK, "option updated successfully", nil)
	}
}

func DeleteOptionHandler(service service.IOptionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid id. the id must be a posetive integer", nil)
		}

		err = service.DeleteOption(model.OptionID(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "error in deleting the option", nil)
		}

		return response.Success(c, fiber.StatusOK, "option deleted", nil)
	}
}
