package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterOptionRoutes registers routes related to option management
func RegisterOptionRoutes(api fiber.Router, s *common.Server) {
	optionRepo := repository.NewOptionRepository(s.DB)
	questionRepo := repository.NewQuestionRepository(s.DB)

	optionService := service.NewOptionService(optionRepo, questionRepo)

	withAuthMiddlewalre := middleware.WithAuthMiddleware(s.DB, s.Cfg.Secret)

	api.Post("", withAuthMiddlewalre, handlers.CreateOptionHandler(optionService))
	api.Get("/:id", withAuthMiddlewalre, handlers.GetOptionHandler(optionService))
	api.Put("/:id", withAuthMiddlewalre, handlers.UpdateOptionHandler(optionService))
	api.Delete("/:id", withAuthMiddlewalre, handlers.DeleteOptionHandler(optionService))
	api.Get("/question/:question_id", withAuthMiddlewalre, handlers.GetOptionsByQuestionIDHandler(optionService))
}
