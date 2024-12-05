package handlers

import (
	"strconv"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func StartHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from context
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		questionnaireIDStr := c.Params("questionnaire_id")
		if questionnaireIDStr == "" {
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id is required", nil)
		}

		// Convert to model.QuestionnaireID and add proper error handling if it's not a valid uint
		questionnaireID, err := strconv.Atoi(questionnaireIDStr)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid questionnaire_id", nil)
		}

		questionResponse, err := service.Start(model.QuestionnaireID(uint(questionnaireID)), model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "questionnaire started", questionResponse)
	}
}

func SubmitHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from context
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		questionIDStr := c.Params("question_id")
		if questionIDStr == "" {
			return response.Error(c, fiber.StatusBadRequest, "question_id is required", nil)
		}

		questionID, err := strconv.Atoi(questionIDStr)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid question_id", nil)
		}

		var ans model.Answer

		if err := c.BodyParser(&ans); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", nil)
		}

		// Call service method
		err = service.Submit(model.QuestionID(uint(questionID)), &ans, model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "answer submitted", nil)
	}
}

func BackHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from context
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		// Call the service method
		questionResponse, err := service.Back(model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "previous question", questionResponse)
	}
}

func NextHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from context
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		// Call the service method
		questionResponse, err := service.Next(model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "next question", questionResponse)
	}
}

func EndHandler(service service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the user from context
		user, ok := c.Locals("user").(model.User)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		err := service.End(model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "questionnaire ended", nil)
	}
}
