package handlers

import (
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/gofiber/fiber/v2"
)

func HelloQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	// closure
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func CreateQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func GetQuestionHandler(service service.IQuestionService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return nil
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
