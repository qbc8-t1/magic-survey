package repository

import (
	"github.com/QBC8-Team1/magic-survey/internal/domain/model"
	"gorm.io/gorm"
)

type IQuestionRepository interface {
	CreateQuestion(question model.Questions) error
	GetQuestionById(ids []uint) (*model.Questions, error)
	UpdateQuestion(question model.Questions) error
	DeleteQuestion(id uint) error
}

type QuestionRepository struct {
	// dependency injection
	db *gorm.DB
}

func NewQuestionRpository(db *gorm.DB) IQuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) CreateQuestion(question model.Questions) error {
	return r.db.Create(&question).Error
}

func (r *QuestionRepository) GetQuestionById(ids []uint) (*model.Questions, error) {
	var question model.Questions
	result := r.db.First(&question, ids)
	return &question, result.Error
}

func (r *QuestionRepository) UpdateQuestion(question model.Questions) error {
	return r.db.Save(&question).Error
}

func (r *QuestionRepository) DeleteQuestion(id uint) error {
	return r.db.Delete(&model.Questions{}, id).Error
}
