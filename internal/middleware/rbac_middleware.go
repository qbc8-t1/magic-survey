package middleware

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CheckBy string

const (
	CHECK_BY_ANSWER_ID        = "answer_id"
	CHECK_BY_QUESTION_ID      = "question_id"
	CHECK_BY_OPTION_ID        = "option_id"
	CHECK_BY_QUESTIONNAIRE_ID = "questionnaire_id"
)

func RegisterRbacMiddlewares(router fiber.Router, db *gorm.DB) {
	rbac := service.NewRbacService(repository.NewRbacRepository(db))

	router.Get("/questionnaires/:questionnaire_id", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_QUESTIONNAIRE_VIEW, "param:questionnaire_id"))
	router.Delete("/questionnaires/:questionnaire_id", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_QUESTIONNAIRE_DELETE, "param:questionnaire_id"))
	router.Put("/questionnaires/:questionnaire_id", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_QUESTIONNAIRE_UPDATE, "param:questionnaire_id"))
	router.Patch("/questionnaires/:questionnaire_id/cancel", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_QUESTIONNAIRE_CANCEL, "param:questionnaire_id"))
	router.Patch("/questionnaires/:questionnaire_id/close", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_QUESTIONNAIRE_CLOSE, "param:questionnaire_id"))

	router.Get("/answers/:id", Gate(CHECK_BY_ANSWER_ID, *rbac, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS, "param:id"))
	router.Post("/answers/", Gate(CHECK_BY_QUESTION_ID, *rbac, model.PERMISSION_ANSWER_CREATE, "body:question_id"))
	router.Put("/answers/:id", Gate(CHECK_BY_ANSWER_ID, *rbac, model.PERMISSION_ANSWER_UPDATE, "param:id"))
	router.Delete("/answers/:id", Gate(CHECK_BY_ANSWER_ID, *rbac, model.PERMISSION_ANSWER_DELETE, "param:id"))

	router.Post("/core/start/:questionnaire_id", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_ANSWER_GIVE, "param:questionnaire_id"))
	router.Post("/submit/:question_id", Gate(CHECK_BY_QUESTION_ID, *rbac, model.PERMISSION_ANSWER_GIVE, "param:question_id"))

	router.Get("/options/:id", Gate(CHECK_BY_OPTION_ID, *rbac, model.PERMISSION_OPTION_VIEW, "param:id"))
	router.Get("/options/question/:question_id", Gate(CHECK_BY_QUESTION_ID, *rbac, model.PERMISSION_OPTION_VIEW, "param:question_id"))
	router.Post("/options/:question_id", Gate(CHECK_BY_QUESTION_ID, *rbac, model.PERMISSION_OPTION_CREATE, "param:question_id"))
	router.Put("/options/:id", Gate(CHECK_BY_OPTION_ID, *rbac, model.PERMISSION_OPTION_UPDATE, "param:id"))
	router.Delete("/options/:id", Gate(CHECK_BY_OPTION_ID, *rbac, model.PERMISSION_OPTION_DELETE, "param:id"))

	router.Post("/questions/", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_QUESTION_CREATE, "body:questionnaire_id"))
	router.Get("/questions/:id", Gate(CHECK_BY_QUESTION_ID, *rbac, model.PERMISSION_QUESTION_VIEW, "param:id"))
	router.Get("/questions/questionnaire/:questionnaire_id", Gate(CHECK_BY_QUESTIONNAIRE_ID, *rbac, model.PERMISSION_QUESTION_VIEW, "param:questionnaire_id"))
	router.Put("/questions/:id", Gate(CHECK_BY_QUESTION_ID, *rbac, model.PERMISSION_QUESTION_UPDATE, "param:id"))
	router.Delete("/questions/:id", Gate(CHECK_BY_QUESTION_ID, *rbac, model.PERMISSION_QUESTION_DELETE, "param:id"))

	router.Post("/superadmin/make-superadmin", SuperadminGate(*rbac, model.PERMISSION_MAKE_NEW_SUPERADMIN))
	router.Post("/superadmin/limit-user-questionnaires-count", SuperadminGate(*rbac, model.PERMISSION_LIMIT_USER_QUESTIONNAIRES_COUNT))
}

func Gate(checkBy CheckBy, rbac service.RbacService, permissionName model.PermissionName, idPlaceName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		localUser := c.Locals("user")
		if localUser == nil {
			// If the user is not set, return unauthorized
			response.Error(c, fiber.StatusUnauthorized, "you are not logged in", nil)
		}

		// Type assert the user to your User struct
		user, ok := localUser.(model.User)
		if !ok {
			// Handle the case where the type assertion fails
			return response.Error(c, fiber.StatusInternalServerError, "failed to get user", nil)
		}

		fieldID, err := getFieldID(c, idPlaceName)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "rbac: failed to retrieve answer id", nil)
		}

		var questionnaire_id model.QuestionnaireID
		switch checkBy {
		case CHECK_BY_QUESTION_ID:
			questionnaire_id, err = rbac.GetQuestionnaireIDByQuestionID(model.QuestionID(fieldID))
		case CHECK_BY_ANSWER_ID:
			questionnaire_id, err = rbac.GetQuestionnaireIDByAnswerID(model.AnswerID(fieldID))
		case CHECK_BY_OPTION_ID:
			questionnaire_id, err = rbac.GetQuestionnaireIDByOptionID(model.OptionID(fieldID))
		case CHECK_BY_QUESTIONNAIRE_ID:
			questionnaire_id = model.QuestionnaireID(fieldID)
		default:
			return response.Error(c, fiber.StatusInternalServerError, "failed to get required data", nil)
		}

		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "failed to get record from db", nil)
		}

		can, err := rbac.CanDo(user.ID, questionnaire_id, permissionName)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong", err.Error())
		}

		if can {
			return c.Next()
		} else {
			return response.Error(c, fiber.StatusForbidden, "you are not allowed to do this", nil)
		}
	}
}

func getFieldID(c *fiber.Ctx, idPlaceName string) (int, error) {
	placeFieldName := strings.Split(idPlaceName, ":")
	if len(placeFieldName) != 2 {
		return 0, errors.New("bad field place name")
	}

	switch placeFieldName[0] {
	case "param":
		return c.ParamsInt(placeFieldName[1])
	case "body":
		body := c.Body()
		var jsonData map[string]interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			return 0, errors.New("field not found in the body")
		}
		value, ok := jsonData[placeFieldName[1]].(float64)
		if !ok {
			return 0, errors.New("field not found in the body")
		}
		return int(value), nil
	default:
		return 0, errors.New("bad field place name")
	}
}

func SuperadminGate(rbac service.RbacService, permissionName model.PermissionName) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			// If the user is not set, return unauthorized
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Type assert the user to your User struct
		currentUser, ok := user.(model.User)
		if !ok {
			// Handle the case where the type assertion fails
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve user")
		}

		can, err := rbac.CanDoAsSuperadmin(currentUser.ID, permissionName)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong", err.Error())
		}

		if can {
			return c.Next()
		} else {
			return response.Error(c, fiber.StatusForbidden, "you are not allowed to do this", nil)
		}
	}
}
