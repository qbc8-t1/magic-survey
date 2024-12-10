package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	logger2 "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type GivePermissionsData struct {
	UserID          uint                     `json:"user_id"`
	QuestionnaireID model.QuestionnaireID    `json:"questionnaire_id"`
	Permissions     []service.PermissionType `json:"permissions"`
}

type RevokePermissionData struct {
	UserID          uint                  `json:"user_id"`
	QuestionnaireID model.QuestionnaireID `json:"questionnaire_id"`
	PermissionName  model.PermissionName  `json:"permission_name"`
}

type HasPermissionData struct {
	QuestionnaireID model.QuestionnaireID `json:"questionnaire_id"`
	PermissionName  model.PermissionName  `json:"permission_name"`
}

func GetAllPermissions(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.JSON(rbacService.GetAllPermissions())
		return nil
	}
}

func GivePermissions(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogRbac))
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		giverUser, ok := localUser.(model.User)
		if !ok {
			logger.Error("something went wrong to get signed in user")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		data := new(GivePermissionsData)
		err := c.BodyParser(data)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "invalid body", err)
		}

		err = rbacService.GivePermissions(uint(giverUser.ID), data.UserID, data.QuestionnaireID, data.Permissions)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to give permissions", nil)
		}

		logger.Info("permissions have given to user")
		return response.Success(c, fiber.StatusCreated, "permissions have given to user", nil)
	}
}

func RevokePermission(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogRbac))

		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("user is not signed in")
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		revokerUser, ok := localUser.(model.User)
		if !ok {
			logger.Error("something went wrong to get signed in user")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		data := new(RevokePermissionData)
		err := c.BodyParser(data)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusUnprocessableEntity, "bad data", nil)
		}

		err = rbacService.RevokePermission(uint(revokerUser.ID), data.UserID, data.QuestionnaireID, data.PermissionName)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "bad data", nil)
		}

		logger.Info("permissions revoked from user")
		return response.Success(c, fiber.StatusOK, "permissions revoked from user", nil)
	}
}

func CanDo(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogRbac))
		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("user is not signed in")
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			logger.Error("something went wrong to get signed in user")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		data := new(HasPermissionData)
		err := c.BodyParser(data)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusUnprocessableEntity, "invalid data entity", nil)
		}

		has, err := rbacService.CanDo(user.ID, data.QuestionnaireID, data.PermissionName)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "internal error", err.Error())
		}

		logger.Info("permission checked")
		return response.Success(c, fiber.StatusOK, "has permission status", has)
	}
}

func GetUser(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogRbac))
		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("user is not signed in")
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			logger.Error("something went wrong to get signed in user")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		user, err := rbacService.GetUser(user.ID)
		if err != nil {
			logger.Error(err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		logger.Info("User is fetched: " + user.Email)
		return c.JSON(model.ToUserResponse(&user))
	}
}

func GetUserRolesWithPermissions(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogRbac))
		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("user is not signed in")
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			logger.Error("something went wrong to get signed in user")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		roles, err := rbacService.GetUserRolesWithPermissions(user.ID)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get roles and permissions", nil)
		}

		if len(roles) == 0 {
			logger.Error("user doesn't have any roles")
			return response.Error(c, fiber.StatusOK, "you don't have any specific roles or permissions", nil)
		}

		logger.Info("roles given")
		return response.Error(c, fiber.StatusOK, "your roles and permissions", roles)
	}
}

func GetUserRolesAndPermissions(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		roles, err := rbacService.GetUserRolesWithPermissions(user.ID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if len(roles) == 0 {
			return c.SendString("user doesn't have any roles")
		}

		return c.JSON(roles)
	}
}

func MakeFakeUser(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := rbacService.MakeFakeUser()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(user)
	}
}

func MakeFakeQuestionnaire(rbacService service.IRbacService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		user, ok := localUser.(model.User)
		if !ok {
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		questionnaire, err := rbacService.MakeFakeQuestionnaire(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(questionnaire)
	}
}
