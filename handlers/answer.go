package handlers

import (
	"strconv"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func HelloAnswerHandler(service service.IAnswerService) fiber.Handler {
	// closure
	return func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "Hello From Answer Handler!", nil)
	}
}

func CreateAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var answerDTO model.CreateAnswerDTO
		if err := c.BodyParser(&answerDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		err := answerDTO.Validate()

		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		err = service.CreateAnswer(&answerDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "Answer created", nil)
	}
}

func GetAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", err)
		}

		res, err := service.GetAnswerByID(uint(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Answer Found", res)
	}
}

func UpdateAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a positive integer", err)
		}

		var answerDTO model.UpdateAnswerDTO
		if err := c.BodyParser(&answerDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		err = answerDTO.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		err = service.UpdateAnswer(uint(id), &answerDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Answer updated successfully", nil)
	}
}

func DeleteAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", err)
		}

		err = service.DeleteAnswer(uint(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Answer Deleted", nil)
	}
}
