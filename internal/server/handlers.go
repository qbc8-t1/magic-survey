package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type jsonResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, jsonResponse{
		Success: true,
		Message: "ok",
	})
}
