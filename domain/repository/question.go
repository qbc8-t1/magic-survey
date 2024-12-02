package domain_repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
)

// IQuestionRepository inteface defines repository methods
type IQuestionRepository interface {
	CreateQuestion(question *model.Question) error
	GetQuestionByID(id model.QuestionID) (*model.Question, error)
	GetQuestionsByQuestionnaireID(questionnaireID model.QuestionnaireID) (*[]model.Question, error)
	UpdateQuestion(question *model.Question) error
	DeleteQuestion(id model.QuestionID) error
}
