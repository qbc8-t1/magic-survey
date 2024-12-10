package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterOptionRoutes registers routes related to option management
func RegisterOptionRoutes(api fiber.Router, s *common.Server) {
	optionRepo := repository.NewOptionRepository(s.DB)
	questionRepo := repository.NewQuestionRepository(s.DB)

	optionService := service.NewOptionService(optionRepo, questionRepo)

	api.Get("/:id", handlers.GetOptionHandler(optionService))
	api.Get("/question/:question_id", handlers.GetOptionsByQuestionIDHandler(optionService))
	api.Post("", handlers.CreateOptionHandler(optionService))
	api.Put("/:id", handlers.UpdateOptionHandler(optionService))
	api.Delete("/:id", handlers.DeleteOptionHandler(optionService))

}
