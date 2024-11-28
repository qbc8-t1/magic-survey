package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	successResp := Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(successResp)
}

func Error(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	errorResp := Response{
		Status:  "error",
		Message: message,
		Error:   err,
	}
	return c.Status(statusCode).JSON(errorResp)
}
