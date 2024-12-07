package middleware

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	jwt2 "github.com/QBC8-Team1/magic-survey/pkg/jwt"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// CustomClaims defines custom JWT claims
type CustomClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

// WithAuthMiddleware validates JWT tokens and fetches user data from the database
func WithAuthMiddleware(db *gorm.DB, secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		logger := GetLogger(c)
		if tokenString == "" {
			logger.Warn("Missing Authorization header", zap.String("path", c.Path()))
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims, err := jwt2.ParseToken(tokenString, []byte(secret))
		if err != nil {
			logger.Warn("Failed to parse token", zap.String("path", c.Path()), zap.Error(err))
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			logger.Warn("Token expired", zap.String("path", c.Path()), zap.Uint("user_id", claims.UserID))
			return response.Error(c, fiber.StatusUnauthorized, "Token has expired", nil)
		}

		var user model.User
		err = db.First(&user, claims.UserID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				logger.Warn("User not found", zap.Uint("user_id", claims.UserID))
				return response.Error(c, fiber.StatusUnauthorized, "User not found", nil)
			}
			logger.Error("Database error", zap.Error(err))
			return response.Error(c, fiber.StatusInternalServerError, "Database error", nil)
		}

		// Store user in context
		c.Locals("user", user)

		return c.Next()
	}
}
