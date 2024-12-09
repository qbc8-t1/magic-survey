package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

func RegisterCoreRoutes(api fiber.Router, s *common.Server) {
	submissionRepo := repository.NewSubmissionRepository(s.DB)
	questionnaireRepo := repository.NewQuestionnaireRepository(s.DB)
	questionRepo := repository.NewQuestionRepository(s.DB)
	optionRepo := repository.NewOptionRepository(s.DB)
	answerRepo := repository.NewAnswerRepository(s.DB)

	coreService := service.NewCoreService(submissionRepo, questionnaireRepo, questionRepo, optionRepo, answerRepo)

	withAuthMiddleware := middleware.WithAuthMiddleware(s.DB, s.Cfg.Secret)

	api.Post("/start/:questionnaire_id", withAuthMiddleware, handlers.StartHandler(coreService))
	api.Post("/submit/:question_id", withAuthMiddleware, handlers.SubmitHandler(coreService))
	api.Post("/back", withAuthMiddleware, handlers.BackHandler(coreService))
	api.Post("/next", withAuthMiddleware, handlers.NextHandler(coreService))
	api.Post("/end", withAuthMiddleware, handlers.EndHandler(coreService))
}
