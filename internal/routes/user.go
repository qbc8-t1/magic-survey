package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes registers routes related to user management
func RegisterUserRoutes(auth fiber.Router, s *common.Server) {
	// ---------------- init requirements
	userRepo := repository.NewUserRepository(s.DB)
	userService := service.NewUserService(userRepo, s.Cfg.Secret, s.Cfg.AuthExpMinute, s.Cfg.AuthRefreshMinute, s.Cfg.Server.MailPass, s.Cfg.Server.FromMail, s.Cfg.Server.MaxSecondForChangeBirthdate)

	// ---------------- middleware
	withAuthMiddleware := middleware.WithAuthMiddleware(s.DB, s.Cfg.Secret)

	// ---------------- routes
	// global
	//auth.Get("user/:id", handlers.ShowUser(*userService))
	auth.Post("signup", handlers.UserCreate(*userService))
	auth.Post("verify", handlers.Verify2FACode(*userService))
	auth.Post("login", handlers.Login(*userService))

	// with auth
	auth.Get("profile", withAuthMiddleware, handlers.ShowProfile(*userService))
	auth.Put("profile", withAuthMiddleware, handlers.UpdateProfile(*userService))
	auth.Post("wallet-balance", withAuthMiddleware, handlers.IncreaseWalletBalance(*userService))
}
