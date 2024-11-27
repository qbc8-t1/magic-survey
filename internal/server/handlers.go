package server

import (
	"github.com/gofiber/fiber/v2"
)

type jsonResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(jsonResponse{
		Success: true,
		Message: "ok",
	})
}
