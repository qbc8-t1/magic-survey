package domain_repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
)

// IQuestionRepository inteface defines repository methods
type IQuestionRepository interface {
	CreateQuestion(question *model.Question) error
	GetQuestionByID(id model.QuestionID) (*model.Question, error)
	GetQuestionsByQuestionnaireID(questionnaireID model.QuestionnaireID) (*[]model.Question, error)
	GetUnansweredQuestions(questionnaireID model.QuestionnaireID, submissionID model.SubmissionID) (*[]model.Question, error)
	UpdateQuestion(question *model.Question) error
	DeleteQuestion(id model.QuestionID) error
	GetQuestionByQuestionIDAndQuestionnaireID(questionID model.QuestionID, questionnaireID model.QuestionnaireID) (*model.Question, error)
}
