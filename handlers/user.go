package handlers

import (
	"errors"
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ShowUser(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return response.Error(c, fiber.StatusBadRequest, "Invalid ID, it must be a positive integer", err)
		}

		res, err := userService.ShowUser(id)

		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "User found", res)
	}
}

func ShowProfile(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// read user from withAuth middleware
		user := c.Locals("user").(model.User)
		res, err := userService.Profile(&user)

		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "User profile found", res)
	}
}

func UpdateProfile(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// base validation
		var dto model.UpdateUserDTO
		if err := c.BodyParser(&dto); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		// all validation params
		user := c.Locals("user").(model.User)
		newUser := model.ToUserModelForUpdate(user, &dto)
		err := newUser.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}

		// call service

		res, err := userService.UpdateUser(&user, &newUser)

		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "User updated successfully", res)
	}
}

func IncreaseCredit(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dto model.IncreaseCreditDTO
		if err := c.BodyParser(&dto); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		value, err := strconv.Atoi(dto.Value)
		if err != nil || value <= 0 || value > 100000000 {
			return response.Error(c, fiber.StatusBadRequest, "Invalid value, it must be a positive integer", err)
		}
		// all validation params
		user := c.Locals("user").(model.User)

		res, err := userService.IncreaseCredit(&user, int64(value))

		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "User updated successfully", res)
	}
}

func UserCreate(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var dto model.CreateUserDTO
		if err := c.BodyParser(&dto); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		user := model.ToUserModel(&dto)
		err := user.Validate()
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid request params", err.Error())
		}
		tokens, err := userService.CreateUser(user)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusCreated, "User Created", tokens)
	}
}

// Verify2FACode handles the verification of the 2FA code
func Verify2FACode(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.Verify2FACodeRequest
		if err := c.BodyParser(&req); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", err)
		}

		tokens, err := userService.Verify2FACode(req.Email, req.Code)
		if err != nil {
			if errors.Is(err, service.ErrInvalid2FACode) {
				return response.Error(c, fiber.StatusUnauthorized, service.ErrInvalid2FACode.Error(), nil)
			}
			return response.Error(c, fiber.StatusInternalServerError, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "2FA verification successful", tokens)
	}
}

func Login(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid request payload", nil)
		}

		tokens, err := userService.LoginUser(&req)

		if errors.Is(err, service.ErrWrongEmailPass) {
			return response.Error(c, fiber.StatusBadRequest, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Login successful", tokens)
	}
}
