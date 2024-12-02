package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

// IQuestionRepository inteface defines repository methods
type IQuestionRepository interface {
	CreateQuestion(question *model.Question) error
	GetQuestionByID(id uint) (*model.Question, error)
	GetQuestionsByID(ids []uint) (*[]model.Question, error)
	UpdateQuestion(question *model.Question) error
	DeleteQuestion(id uint) error
}
