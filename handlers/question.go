package handlers

import (
	"strconv"

	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func HelloQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	// closure
	return func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "Hello From Question Handler", nil)
	}
}

func CreateQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func GetQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", err)
		}

		res, err := service.GetQuestionByID(uint(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Question Found", res)
	}
}

func GetQuestionsHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func UpdateQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func DeleteQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
