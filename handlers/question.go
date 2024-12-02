package handlers

import (
	"strconv"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func HelloQuestionHandler(service service.IQuestionService) fiber.Handler {
	// closure
	return func(c *fiber.Ctx) error {
		return response.Success(c, fiber.StatusOK, "Hello From Question Handler!", nil)
	}
}

func CreateQuestionHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var questionDTO model.CreateQuestionDTO
		if err := c.BodyParser(&questionDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		err := questionDTO.Validate()

		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		err = service.CreateQuestion(&questionDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "Question created", nil)
	}
}

func GetQuestionHandler(service service.IQuestionService) fiber.Handler {
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

func GetQuestionsByQuestionnaireIDHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		questionnaireIdStr := c.Params("questionnaire_id")
		questionnaireId, err := strconv.Atoi(questionnaireIdStr)
		if err != nil || questionnaireId <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", err)
		}

		questions, err := service.GetQuestionsByQuestionnaireID(uint(questionnaireId))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to fetch questions", err)
		}

		return response.Success(c, fiber.StatusOK, "Questions retrieved successfully", questions)
	}
}

func UpdateQuestionHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a positive integer", err)
		}

		var questionDTO model.UpdateQuestionDTO
		if err := c.BodyParser(&questionDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		err = questionDTO.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		err = service.UpdateQuestion(uint(id), &questionDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Question updated successfully", nil)
	}
}

func DeleteQuestionHandler(service service.IQuestionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", err)
		}

		err = service.DeleteQuestion(uint(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Question Deleted", nil)
	}
}
