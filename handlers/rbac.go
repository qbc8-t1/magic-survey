package handlers

import (
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type GivePermissionsData struct {
	UserID          uint                     `json:"user_id"`
	QuestionnaireID uint                     `json:"questionnaire_id"`
	Permissions     []service.PermissionType `json:"permissions"`
}

type RevokePermissionData struct {
	UserID          uint   `json:"user_id"`
	QuestionnaireID uint   `json:"questionnaire_id"`
	PermissionName  string `json:"permission_name"`
}

type HasPermissionData struct {
	UserID          uint   `json:"user_id"`
	QuestionnaireID uint   `json:"questionnaire_id"`
	PermissionName  string `json:"permission_name"`
}

type GetUsersWithVisibleAnswersData struct {
	UserID          uint `json:"user_id"`
	QuestionnaireID uint `json:"questionnaire_id"`
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
			return response.Error(c, fiber.StatusBadRequest, "user id param invalid", err)
		}

		data := new(GivePermissionsData)
		err = c.BodyParser(data)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		err = rbacService.GivePermissions(uint(giverUserId), data.UserID, data.QuestionnaireID, data.Permissions)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to process", err)
		}

		return response.Success(c, fiber.StatusCreated, "permissions gived to user", nil)
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

		if len(roles) == 0 {
			return c.SendString("user doesn't have any roles")
		}

		return c.JSON(roles)
	}
}

func RevokePermission(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		revokerUserID, err := c.ParamsInt("userid")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		data := new(RevokePermissionData)
		err = c.BodyParser(data)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		err = rbacService.RevokePermission(uint(revokerUserID), data.UserID, data.QuestionnaireID, data.PermissionName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		c.Status(fiber.StatusOK).SendString("permissions revoked from user")
		return nil
	}
}

func CanDo(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := new(HasPermissionData)
		err := c.BodyParser(data)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		has, err := rbacService.CanDo(data.UserID, data.QuestionnaireID, data.PermissionName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(has)
	}
}

func GetUsersWithVisibleAnswers(rbacService service.RbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := new(GetUsersWithVisibleAnswersData)
		err := c.BodyParser(&data)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		usersIDs, err := rbacService.GetUsersWithVisibleAnswers(data.QuestionnaireID, data.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(usersIDs)
	}
}
