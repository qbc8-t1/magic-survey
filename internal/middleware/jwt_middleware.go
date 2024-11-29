package middleware

import (
	"fmt"
	"github.com/QBC8-Team1/magic-survey/domain/model"
	jwt2 "github.com/QBC8-Team1/magic-survey/pkg/jwt"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
		// Extract token from Authorization header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		fmt.Println(tokenString)

		claims, err := jwt2.ParseToken(tokenString, []byte(secret))
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, err.Error(), nil)
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			return response.Error(c, fiber.StatusUnauthorized, "Token has expired", nil)
		}

		var user model.User

		err = db.First(&user, claims.UserID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return response.Error(c, fiber.StatusUnauthorized, "", nil)
			}
			return response.Error(c, fiber.StatusInternalServerError, "", nil)
		}

		// Store user in context
		c.Locals("user", user)

		return c.Next()
	}
}
