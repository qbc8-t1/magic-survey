package server

import (
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type jsonResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func healthCheck(c *fiber.Ctx) error {
	user := map[string]interface{}{
		"data": 222,
	}

	return response.Success(c, 201, "Healthy", user)

}
