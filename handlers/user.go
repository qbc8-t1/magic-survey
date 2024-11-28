package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/domain/repository"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// 1. +get data
// 2. +init db connection
// 3. password hashing (add salt row in db)
// 4. signup
// 5. login
// 6. setup SMTP server for sending email
// 7. send verification code
func UserCreate(repo repository.IUserRepository) func(c *fiber.Ctx) error {
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

		if err := repo.CreateUser(userModel); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create user",
			})
		}

		userResponse := model.ToUserResponse(userModel)

		return response.Success(c, 201, "userCreated!!", userResponse)
	}
}
