package handlers

import (
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, 200, "Healthy", nil)
}
