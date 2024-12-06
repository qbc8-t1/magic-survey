package middleware

import (
	"github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggingMiddleware sets up a basic logger for each request
func WithLoggingMiddleware(appLogger *logger.AppLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceId := uuid.NewString()
		requestLogger := appLogger.WithFields(map[string]interface{}{
			"path":      c.Path(),
			"method":    c.Method(),
			"client_ip": c.IP(),
			"trace_id":  traceId,
		})

		c.Locals("logger", requestLogger)

		return c.Next()
	}
}

// GetLogger retrieves the logger from the Fiber context
func GetLogger(c *fiber.Ctx) *zap.Logger {
	logger, ok := c.Locals("logger").(*zap.Logger)
	if !ok {
		return zap.L() // Fallback to the global logger
	}
	return logger
}
