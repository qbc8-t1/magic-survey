package handlers

import (
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/gofiber/fiber/v2"
)

func HelloHandlerQuestion(service service.IQuestionService) func(c *fiber.Ctx) error {
	// closure
	return func(c *fiber.Ctx) error {
		return nil
	}
}
