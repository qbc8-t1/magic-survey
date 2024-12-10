package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterNotificationRoutes registers routes related to notification management
func RegisterNotificationRoutes(auth fiber.Router, s *common.Server) {
	// Initialize repositories and services
	notificationRepo := repository.NewNotificationRepository(s.DB)
	notificationService := service.NewNotificationService(notificationRepo)
	withAuthMiddleware := middleware.WithAuthMiddleware(s.DB, s.Cfg.Secret)

	// Routes with authentication
	auth.Post("/", withAuthMiddleware, handlers.CreateNotification(*notificationService))
	//auth.Put("notification/:id/seen", withAuthMiddleware, handlers.MarkNotificationAsSeen(*notificationService))
	//auth.Get("notifications", withAuthMiddleware, handlers.ListNotifications(*notificationService))
}
