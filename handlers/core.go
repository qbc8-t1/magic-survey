package handlers

import (
	"errors"
	"strconv"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// StartHandler starts a questionnaire
func StartHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(model.User)
		qidStr := c.Params("questionnaire_id")
		qid, err := strconv.ParseUint(qidStr, 10, 64)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid questionnaire_id", nil)
		}
		question, err := svc.Start(model.QuestionnaireID(qid), model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return response.Success(c, fiber.StatusOK, "started", question)
	}
}

func SubmitHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract user from context
		user, ok := c.Locals("user").(model.User)
		if !ok || user.ID <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "userID is required and must be greater than 0", nil)
		}

		// Parse questionID from params
		qidStr := c.Params("question_id")
		qid, err := strconv.ParseUint(qidStr, 10, 64)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid question_id", nil)
		}

		// Parse the answer DTO from the request body
		var dto model.CreateAnswerDTO
		if err := c.BodyParser(&dto); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
		}

		// Validate the DTO
		err = dto.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		// Map DTO to Answer model
		answer := model.ToAnswerModel(&dto)

		// Call the service to submit the answer
		err = svc.Submit(model.QuestionID(qid), answer, model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "answer submitted", nil)
	}
}

// BackHandler moves user to the previous question
func BackHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(model.User)
		q, err := svc.Back(model.UserId(user.ID))
		if err != nil {
			if errors.Is(err, service.ErrCannotGoBack) {
				return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
			}
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return response.Success(c, fiber.StatusOK, "previous question", q)
	}
}

// NextHandler moves user to the next question
func NextHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(model.User)
		q, err := svc.Next(model.UserId(user.ID))
		if err != nil {
			if errors.Is(err, service.ErrNoNextQuestion) {
				return response.Error(c, fiber.StatusNotFound, "no more questions left", nil)
			}
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return response.Success(c, fiber.StatusOK, "next question", q)
	}
}

// EndHandler finalizes the questionnaire submission
func EndHandler(svc service.ICoreService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(model.User)
		err := svc.End(model.UserId(user.ID))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return response.Success(c, fiber.StatusOK, "questionnaire ended", nil)
	}
}
