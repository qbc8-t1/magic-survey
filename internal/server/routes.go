package server

import (
	repository "github.com/QBC8-Team1/magic-survey/domain/infra"
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, s *Server) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	api := app.Group("/api")
	auth := api.Group("/v1/auth")

	// 1. +get data
	// 2. +init db connection
	// 3. password hashing(add salt row in db)
	// 4. signup
	// 5.login
	// 6. setup smtp server for sending email
	// 7. send verification code
	auth.Post("signup", func(c *fiber.Ctx) error {
		var user model.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		userRepo := repository.NewUserRepository(s.db)

		userRepo.CreateUser(user)

		return c.JSON(fiber.Map{
			"status": "ok",
			"data":   user,
		})
	})

}
