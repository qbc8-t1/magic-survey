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

type GetAnswerData struct {
	QuestionID uint `json:"question_id"`
	UserID     uint `json:"user_id"`
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

		// get questionnaire id
		questionnaireID, err := c.ParamsInt("questionnaire_id")
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "bad questionnaire id", err.Error())
		}

		// get questionnaire object by id
		questionnaire, err := questionnaireService.GetQuestionnaireByID(model.QuestionnaireID(uint(questionnaireID)))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get questionnaireID", err.Error())
		}

		// parse body of request
		data := new(GetAnswerData)
		err = c.BodyParser(&data)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusUnprocessableEntity, "faild to parse body", err.Error())
		}

		// check if question is for the questionnaire
		is, err := questionService.IsQuestionForQuestionnaire(data.QuestionID, questionnaire.ID)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get question", err.Error())
		}
		if !is {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusBadRequest, "question is not for questionnaire", nil)
		}

		// user can see his/her answer
		if data.UserID == loggedInUser.ID {
			answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserId(data.UserID))
			if err != nil {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", err.Error())
			}
			return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))
		}

		// switch on questionnaire visibility type
		switch questionnaire.AnswersVisibleFor {

		// if everybody can see it
		case model.QuestionnaireVisibilityEverybody:
			// check if logged in user is superadmin or the owner of the questionnaire
			isSuperadmin, _ := rbacService.CanDoAsSuperadmin(loggedInUser.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)
			if questionnaire.OwnerID == loggedInUser.ID || isSuperadmin {
				answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserId(data.UserID))
				if err != nil {
					logger.Error(err.Error())
					return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", err.Error())
				}

				logger.Info("user answer was successful")
				return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))
			}

			// check if user has permission
			has, err := rbacService.HasPermission(loggedInUser.ID, questionnaire.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)
			if err != nil {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusInternalServerError, "failed to check permission", err.Error())
			}
			if !has {
				logger.Error("you don't have permission to see the answer")
				return response.Error(c, fiber.StatusForbidden, "you don't have permission to see the answer", nil)
			}

			// get selected users which users can see their answers
			usersIDs, err := rbacService.GetUsersIDsWithVisibleAnswers(questionnaire.ID, loggedInUser.ID)
			if err != nil {
				logger.Error(err.Error())
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
				return response.Error(c, fiber.StatusForbidden, msg, nil)
			}
			answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserId(data.UserID))
			if err != nil {
				logger.Error(err.Error())
				return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", err.Error())
			}

			logger.Info("user answer")
			return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))

		case model.QuestionnaireVisibilityAdminAndOwner:
			isSuperadmin, _ := rbacService.CanDoAsSuperadmin(loggedInUser.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)

			if questionnaire.OwnerID == loggedInUser.ID || isSuperadmin {
				answers, err := answerService.GetUserAnswers(model.QuestionID(data.QuestionID), model.UserId(data.UserID))
				if err != nil {
					logger.Error(err.Error())
					return response.Error(c, fiber.StatusInternalServerError, "failed to get answer", err.Error())
				}
				logger.Info("get another user answer successful")
				return response.Success(c, fiber.StatusOK, "user answer", model.ToAnswerSummaryResponses(answers))
			}

			logger.Error("you don't have permission to see the answer")
			return response.Error(c, fiber.StatusForbidden, "you don't have permission to see the answer", nil)
		case model.QuestionnaireVisibilityNobody:
			logger.Error("answer is not visible for anyone")
			return response.Success(c, fiber.StatusOK, "answer is not visible for anyone", nil)
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
			return response.Error(c, fiber.StatusBadRequest, "bad questionnaire id", err.Error())
		}

		// get questionnaire object by id
		questionnaire, err := questionnaireService.GetQuestionnaireByID(model.QuestionnaireID(uint(questionnaireID)))
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get questionnaireID", err.Error())
		}

		has, err := rbacService.HasPermission(loggedInUser.ID, questionnaire.ID, model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to check permission", err.Error())
		}
		if !has {
			logger.Error("you don't have specific permission to see the answer")
			return response.Error(c, fiber.StatusForbidden, "you don't have specific permission to see the answer", nil)
		}

		usersIDs, err := rbacService.GetUsersIDsWithVisibleAnswers(questionnaire.ID, loggedInUser.ID)
		if err != nil {
			logger.Error(err.Error())
			return response.Error(c, fiber.StatusInternalServerError, "failed to get users ids with visible answers", err.Error())
		}

		logger.Info("you have permission to see users answers")
		return response.Success(c, fiber.StatusOK, "you have specific permission to see these users answer for selected questionnaire", usersIDs)
	}
}
