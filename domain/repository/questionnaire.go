package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type IQuestionnaireRepository interface {
	GetQuestionnaireByID(questionnaireID uint) (model.Questionnaire, error)
}