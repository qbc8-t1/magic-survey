package domain_repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/google/uuid"
)

// IQuestionRepository inteface defines repository methods
type IQuestionRepository interface {
	CreateQuestion(question *model.Question) error
	GetQuestionByID(id uuid.UUID) (*model.Question, error)
	GetQuestionsByID(ids []uuid.UUID) (*[]model.Question, error)
	UpdateQuestion(question *model.Question) error
	DeleteQuestion(id uuid.UUID) error
}
