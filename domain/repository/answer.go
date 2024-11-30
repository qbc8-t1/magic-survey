package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

// IAnswerRepository interface defines the repository methods
type IAnswerRepository interface {
	CreateAnswer(answer *model.Answer) error
	GetAnswerByID(id uint) (*model.Answer, error)
	UpdateAnswer(answer *model.Answer) error
	DeleteAnswer(id uint) error
}
