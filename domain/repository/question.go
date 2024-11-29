package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type IQuestionRepository interface {
	CreateQuestion(question model.Question) error
	GetQuestionById(ids []uint) (*model.Question, error)
	UpdateQuestion(question model.Question) error
	DeleteQuestion(id uint) error
}
