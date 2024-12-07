package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type IQuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *model.Questionnaire) (model.Questionnaire, error)
	GetQuestionnaireByID(questionnaireID uint) (model.Questionnaire, error)
	UpdateQuestionaire(questionnaire *model.Questionnaire) error
	DeleteQuestionnaire(id uint) error
	GetUserQuestionnairesCount(userID uint) (int64, error)
}
