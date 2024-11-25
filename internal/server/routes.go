package server

import (
	"github.com/labstack/echo/v4"
)

func registerRoutes(e *echo.Echo, s *Server) {
	s.router.GET("/health", healthCheck)

}
