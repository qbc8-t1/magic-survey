package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 1. +get data
// 2. +init db connection
// 3. password hashing (add salt row in db)
// 4. signup
// 5. login
// 6. setup SMTP server for sending email
// 7. send verification code
func UserCreate(db *gorm.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var userDTO model.CreateUserDTO
		if err := c.BodyParser(&userDTO); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		userDTO.Password = string(hashedPassword)

		userModel := model.ToUserModel(&userDTO)

		userRepo := repository.NewUserRepository(db)
		if err := userRepo.CreateUser(userModel); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create user",
			})
		}

		userResponse := model.ToUserResponse(userModel)

		return response.Success(c, 201, "userCreated!!", userResponse)
	}
}
