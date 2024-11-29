package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRpository(db *gorm.DB) domain_repository.IQuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) CreateQuestion(question model.Question) error {
	return r.db.Create(&question).Error
}

func (r *QuestionRepository) GetQuestionById(ids []uint) (*model.Question, error) {
	var question model.Question
	result := r.db.First(&question, ids)
	return &question, result.Error
}

func (r *QuestionRepository) UpdateQuestion(question model.Question) error {
	return r.db.Save(&question).Error
}

func (r *QuestionRepository) DeleteQuestion(id uint) error {
	return r.db.Delete(&model.Question{}, id).Error
}
