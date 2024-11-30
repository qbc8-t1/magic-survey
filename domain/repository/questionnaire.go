package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type IQuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *model.Questionnaire) error
	GetQuestionnaireByID(id uint) (*model.Questionnaire, error)
	UpdateQuestionare(questionnaire *model.Questionnaire) error
	DeleteQuestionnaire(id uint) error
}
