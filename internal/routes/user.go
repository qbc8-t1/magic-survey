package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// registerUserRoutes registers routes related to user management
func RegisterUserRoutes(auth fiber.Router, s *common.Server) {
	userRepo := repository.NewUserRepository(s.DB)

	auth.Post("signup", handlers.UserCreate(userRepo))
}
