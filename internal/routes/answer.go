package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterAnswerRoutes registers routes related to answer management
func RegisterAnswerRoutes(app fiber.Router, s *common.Server) {
	answerRepo := repository.NewAnswerRepository(s.DB)
	answerService := service.NewAnswerService(answerRepo)

	app.Get("/hello", handlers.HelloHandlerAnswer(answerService))

}
