package handlers

import (
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/gofiber/fiber/v2"
)

func HelloHandlerAnswer(service service.IAnswerService) func(ctx *fiber.Ctx) error {
	// closure
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
