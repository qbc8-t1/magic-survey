package middleware

import (
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/gofiber/fiber/v2"
	"time"
)

func WithLogger(s *common.Server) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		log := s.Logger

		log.Infof("Request: %s %s from %s", c.Method(), c.Path(), c.IP())

		err := c.Next()

		log.Infof(
			"Response: Status %d - Duration %v - Path %s",
			c.Response().StatusCode(),
			time.Since(start),
			c.Path(),
		)

		return err
	}
}
