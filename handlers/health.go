package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	return response.Success(c, 200, "Hello "+user.GetFullName(), nil)
}
