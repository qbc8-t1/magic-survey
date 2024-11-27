package response

import "github.com/gofiber/fiber/v2"

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	successResp := SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(successResp)
}

func Error(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	errorResp := ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   err,
	}
	return c.Status(statusCode).JSON(errorResp)
}
