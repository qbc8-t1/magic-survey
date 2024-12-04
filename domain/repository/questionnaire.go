package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type IQuestionnaireRepository interface {
	GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (*model.Questionnaire, error)
	GetFirstQuestion(questionnaireID model.QuestionnaireID) (*model.Question, error)
	GetNextQuestion(questionnaireID model.QuestionnaireID, currentOrder int) (*model.Question, error)
	GetPreviousQuestion(questionnaireID model.QuestionnaireID, currentOrder int) (*model.Question, error)
}
