package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes registers routes related to user management
func RegisterUserRoutes(auth fiber.Router, s *common.Server) {
	userRepo := repository.NewUserRepository(s.DB)
	userService := service.NewUserService(userRepo, s.Cfg.Secret, s.Cfg.AuthExpMinute, s.Cfg.AuthRefreshMinute, s.Cfg.Server.SenderMail, s.Cfg.Server.MailPass)

	auth.Post("signup", handlers.UserCreate(*userService))
	auth.Post("verify", handlers.Verify2FACode(*userService))
	auth.Post("login", handlers.Login(*userService))
}
