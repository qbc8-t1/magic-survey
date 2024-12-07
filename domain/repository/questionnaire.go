package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type IQuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *model.Questionnaire) (model.Questionnaire, error)
	GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (model.Questionnaire, error)
	UpdateQuestionaire(questionnaireID model.QuestionnaireID, questionnaire *model.Questionnaire) error
	DeleteQuestionnaire(questionnaireID model.QuestionnaireID) error
	GetUserQuestionnairesCount(userID model.UserId) (int64, error)
}
