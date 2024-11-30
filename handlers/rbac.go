package handlers

import (
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/gofiber/fiber/v2"
)

type GivePermissionsBody struct {
	UserID          uint                     `json:"user_id"`
	QuestionnaireID uint                     `json:"questionnaire_id"`
	Permissions     []service.PermissionType `json:"permissions"`
}

func GetAllPermissions(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.JSON(rbacService.GetAllPermissions())
		return nil
	}
}

func GivePermissions(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		giverUserId, err := c.ParamsInt("userid")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		data := GivePermissionsBody{}
		err = c.BodyParser(&data)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		err = rbacService.GivePermissions(uint(giverUserId), data.UserID, data.QuestionnaireID, data.Permissions)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		c.Status(fiber.StatusOK).SendString("permissions gived to user")
		return nil
	}
}

func MakeFakeUser(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := rbacService.MakeFakeUser()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(user)
	}
}

func MakeFakeQuestionnaire(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := c.ParamsInt("userid")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		questionnaire, err := rbacService.MakeFakeQuestionnaire(uint(userID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(questionnaire)
	}
}

func GetUser(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := c.ParamsInt("userid")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		user, err := rbacService.GetUser(uint(userID))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		c.JSON(user)
		return nil
	}
}

func GetUserRolesWithPermissions(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := c.ParamsInt("userid")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		roles, err := rbacService.GetUserRolesWithPermissions(uint(userID))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		c.JSON(roles)
		return nil
	}
}

// func RevokePermission(rbacService service.RbacService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		return nil
// 	}
// }

// func HasPermission(rbacService service.RbacService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		return nil
// 	}
// }
