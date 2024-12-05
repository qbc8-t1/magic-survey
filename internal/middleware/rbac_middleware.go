package middleware

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRbacMiddlewares(router fiber.Router, db *gorm.DB) {
	rbac := service.NewRbacService(repository.NewRbacRepository(db))

	router.Get("/questionnaires/:questionnaire_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTIONNAIRE_VIEW))
	router.Post("/questionnaires/:questionnaire_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTIONNAIRE_EDIT))

	router.Post("/questionnaires/:questionnaire_id/questions", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_CREATE))
	router.Get("/questionnaires/:questionnaire_id/questions/:question_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_VIEW))
	router.Put("/questionnaires/:questionnaire_id/questions/:question_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_UPDATE))
	router.Delete("/questionnaires/:questionnaire_id/questions/:question_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_DELETE))

	router.Post("/questionnaires/:questionnaire_id/options", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_OPTION_CREATE))
	router.Get("/questionnaires/:questionnaire_id/options/:option_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_OPTION_VIEW))
	router.Put("/questionnaires/:questionnaire_id/options/:option_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_OPTION_UPDATE))
	router.Use("/questionnaires/:questionnaire_id/options/:option_id", QuestionnaireGate(*rbac, model.PERMISSION_QUESTION_OPTION_DELETE))

	router.Post("/questionnaires/:questionnaire_id/answers", QuestionnaireGate(*rbac, model.PERMISSION_ANSWER_GIVE))
	router.Put("/questionnaires/:questionnaire_id/answers/:answer_id", QuestionnaireGate(*rbac, model.PERMISSION_ANSWER_UPDATE))
	router.Delete("/questionnaires/:questionnaire_id/answers/:answer_id", QuestionnaireGate(*rbac, model.PERMISSION_ANSWER_GIVEBACK))

	router.Get("/questionnaires/:questionnaire_id/reports", QuestionnaireGate(*rbac, model.PERMISSION_REPORTS_VIEW))

	router.Post("/questionnaires/:questionnaire_id/give-permissions", QuestionnaireGate(*rbac, model.PERMISSION_GIVE_OR_TAKE_PERMISSION))
	router.Post("/questionnaires/:questionnaire_id/revoke-permissions", QuestionnaireGate(*rbac, model.PERMISSION_GIVE_OR_TAKE_PERMISSION))

	// router.Post("/superadmin/limit-user-questionnaires-count", SuperadminGate(*rbac, model.PERMISSION_LIMIT_USERS_QUESTIONNAIRES_COUNT))
	router.Post("/superadmin/make-superadmin", SuperadminGate(*rbac, model.PERMISSION_MAKE_NEW_SUPERADMIN))

	// router.Get("/questionnaires/:questinnaire_id/answers/:answer_id/user/:user_id")
	// PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS = "questionnaire_see_selected_users_answers"
}

func QuestionnaireGate(rbac service.RbacService, permissionName string) fiber.Handler {
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

		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "questionnaire_id is required", err.Error())
		}

		can, err := rbac.CanDo(currentUser.ID, uint(questionnaireID), permissionName)
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

func SuperadminGate(rbac service.RbacService, permissionName string) fiber.Handler {
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
