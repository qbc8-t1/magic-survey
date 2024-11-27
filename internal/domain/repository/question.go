package repository

import (
	"github.com/QBC8-Team1/magic-survey/internal/domain/model"
	"gorm.io/gorm"
)

type IQuestionRepository interface {
	Add(question model.Questions) (model.Questions, error)
	Get(ids []uint) ([]model.Questions, error)
	Update(question model.Questions) error
	Delete(id uint) error
}

type QuestionRepository struct {
	// dependency injection
	db *gorm.DB
}

func NewQuestionRpository(db *gorm.DB) IQuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Add(question model.Questions) (model.Questions, error) {
	result := r.db.Create(&question)
	return question, result.Error
}

func (r *QuestionRepository) Get(ids []uint) ([]model.Questions, error) {
	var questions []model.Questions
	result := r.db.Find(&questions, ids)
	return questions, result.Error
}

func (r *QuestionRepository) Update(question model.Questions) error {
	result := r.db.Save(&question)
	return result.Error
}

func (r *QuestionRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Questions{}, id)
	return result.Error
}
