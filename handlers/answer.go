package handlers

import (
	"strconv"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func CreateAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var answerDTO model.CreateAnswerDTO
		if err := c.BodyParser(&answerDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err := answerDTO.Validate()

		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", nil)
		}

		err = service.CreateAnswer(&answerDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "error in creating the answer", nil)
		}

		return response.Success(c, fiber.StatusCreated, "answer created", nil)
	}
}

func GetAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a posetive integer", nil)
		}

		res, err := service.GetAnswerByID(model.AnswerID(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "error in retrieving the answer", nil)
		}

		return response.Success(c, fiber.StatusOK, "answer Found", res)
	}
}

func UpdateAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid ID. the ID must be a positive integer", nil)
		}

		var answerDTO model.UpdateAnswerDTO
		if err := c.BodyParser(&answerDTO); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		err = answerDTO.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", nil)
		}

		err = service.UpdateAnswer(model.AnswerID(id), &answerDTO)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "error in updating the answer", nil)
		}

		return response.Success(c, fiber.StatusOK, "answer updated successfully", nil)
	}
}

func DeleteAnswerHandler(service service.IAnswerService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "invalid id. the id must be a posetive integer", nil)
		}

		err = service.DeleteAnswer(model.AnswerID(id))
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "error in deleting the answer", nil)
		}

		return response.Success(c, fiber.StatusOK, "answer deleted", nil)
	}
}
