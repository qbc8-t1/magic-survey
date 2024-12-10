package handlers

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	logger2 "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GetAnswerData struct {
	QuestionID model.QuestionID `json:"question_id"`
	UserID     uint             `json:"user_id"`
}

func GetAnotherUserAnswer(answerService service.IAnswerService, rbacService service.IRbacService, questionService service.IQuestionService, questionnaireService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogVisibleAnswer))
		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("unauthorized")
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		// cast to model user
		loggedInUser, ok := localUser.(model.User)
		if !ok {
			logger.Error("unauthorized")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		// parse body of request
		data := new(GetAnswerData)
		err := c.BodyParser(&data)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusUnprocessableEntity, "faild to parse body", nil)
		}

		// check if question is for the questionnaire
		question, err := questionService.GetQuestionByID(data.QuestionID)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get question", nil)
		}
		questionnaire, err := questionnaireService.GetQuestionnaireByID(question.QuestionnaireID)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get questionnaire", nil)
		}

		// user can see his/her answer
		if data.UserID == loggedInUser.ID {
			answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserID(data.UserID))
			if err != nil {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", nil)
			}
			return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))
		}

		// switch on questionnaire visibility type
		switch questionnaire.AnswersVisibleFor {

		// if everybody can see it
		case model.QuestionnaireVisibilityEverybody:
			answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserID(data.UserID))
			if err != nil {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", nil)
			}

			logger.Info("user answer")
			return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))
		case model.QuestionnaireVisibilityAdminAndOwner:
			isSuperadmin, _ := rbacService.CanDoAsSuperadmin(loggedInUser.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)

			if questionnaire.OwnerID == loggedInUser.ID || isSuperadmin {
				answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserID(data.UserID))
				if err != nil {
					logger.Error(err.Error())
					return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", nil)
				}
				logger.Info("get another user answer successful")
				return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))
			}

			// check if user has permission
			has, err := rbacService.HasPermission(loggedInUser.ID, questionnaire.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)
			if err != nil {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusInternalServerError, "failed to check permission", nil)
			}
			if !has {
				logger.Error("you don't have permission to see the answer")
				return response.Error(c, fiber.StatusForbidden, "you don't have permission to see the answer", nil)
			}

			// get selected users which users can see their answers
			usersIDs, err := rbacService.GetUsersIDsWithVisibleAnswers(questionnaire.ID, loggedInUser.ID)
			if err != nil {
				logger.Error(err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					msg := "you don't have permission to see this user answer for this questionnaire"
					logger.Error(msg)
					return response.Error(c, fiber.StatusForbidden, "you don't have permission to see the answer", nil)
				}
				return response.Error(c, fiber.StatusInternalServerError, "failed to get users ids with visible answers", err.Error())
			}
			found := false
			for _, userID := range usersIDs {
				if data.UserID == userID {
					found = true
				}
			}
			if !found {
				msg := "you don't have permission to see this user answer for this questionnaire"
				logger.Error(msg)
				return response.Error(c, fiber.StatusForbidden, "you don't have permission to see the answer", nil)
			}
			answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserID(data.UserID))
			if err != nil {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", nil)
			}

			logger.Info("user answer")
			return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))
		case model.QuestionnaireVisibilityNobody:
			logger.Error("answer is not visible for anyone")
			return response.Error(c, fiber.StatusForbidden, "answer of this questionnaire is not visible for anyone", nil)
		}
		return nil
	}
}

func GetUsersWithVisibleAnswers(rbacService service.IRbacService, questionnaireService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogVisibleAnswer))

		// get user from local
		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("unauthorized")
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		// cast to model user
		loggedInUser, ok := localUser.(model.User)
		if !ok {
			logger.Error("unauthorized")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		// get questionnaire id
		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "bad questionnaire id", nil)
		}

		// get questionnaire object by id
		questionnaire, err := questionnaireService.GetQuestionnaireByID(model.QuestionnaireID(questionnaireID))

		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get questionnaireID", nil)
		}

		has, err := rbacService.HasPermission(loggedInUser.ID, questionnaire.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to check permission", nil)
		}
		if !has {
			logger.Error("you don't have specific permission to see the answer")
			return response.Error(c, fiber.StatusForbidden, "you don't have specific permission to see the answer", nil)
		}

		usersIDs, err := rbacService.GetUsersIDsWithVisibleAnswers(questionnaire.ID, loggedInUser.ID)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get users ids with visible answers", nil)
		}

		logger.Info("you have permission to see users answers")
		return response.Success(c, fiber.StatusOK, "you have specific permission to see these users answer for selected questionnaire", usersIDs)
	}
}

func GetAnswers(rbacService service.IRbacService, questionnaireService service.IQuestionnaireService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := middleware.GetLogger(c).With(zap.String("category", logger2.LogVisibleAnswer))

		// get user from local
		localUser := c.Locals("user")
		if localUser == nil {
			logger.Error("unauthorized")
			return response.Error(c, fiber.StatusUnauthorized, "user is not signed in", nil)
		}

		// cast to model user
		loggedInUser, ok := localUser.(model.User)
		if !ok {
			logger.Error("unauthorized")
			return response.Error(c, fiber.StatusInternalServerError, "something went wrong to get signed in user", nil)
		}

		// get questionnaire id
		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "bad questionnaire id", nil)
		}

		// get questionnaire object by id
		questionnaire, err := questionnaireService.GetQuestionnaireByID(model.QuestionnaireID(questionnaireID))

		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get questionnaireID", nil)
		}

		has, err := rbacService.CanDo(loggedInUser.ID, questionnaire.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to check permission", nil)
		}
		if !has {
			logger.Error("you don't have specific permission to see the answer")
			return response.Error(c, fiber.StatusForbidden, "you don't have specific permission to see the answer", nil)
		}

		result := rbacService.GetQuestionnaireAnswers(questionnaire.ID)
		return response.Success(c, fiber.StatusOK, "answers", result)

	}
}
