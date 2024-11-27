package repository

import (
	"github.com/QBC8-Team1/magic-survey/internal/domain/model"
	"gorm.io/gorm"
)

type IQuestionnaireRepository interface {
	Add(questionnaire model.Questionnaires) (model.Questionnaires, error)
	Get(ids []uint) ([]model.Questionnaires, error)
	Update(questionnaire model.Questionnaires) error
	Delete(id uint) error
}

type QuestionnaireRepository struct {
	// dependency injection
	db *gorm.DB
}

func NewQuestionnaireRpository(db *gorm.DB) IQuestionnaireRepository {
	return &QuestionnaireRepository{db: db}
}

func (r *QuestionnaireRepository) Add(questionnaire model.Questionnaires) (model.Questionnaires, error) {
	result := r.db.Create(&questionnaire)
	return questionnaire, result.Error
}

func (r *QuestionnaireRepository) Get(ids []uint) ([]model.Questionnaires, error) {
	var questionnaires []model.Questionnaires
	result := r.db.Find(&questionnaires, ids)
	return questionnaires, result.Error
}

func (r *QuestionnaireRepository) Update(questionnaire model.Questionnaires) error {
	result := r.db.Save(&questionnaire)
	return result.Error
}

func (r *QuestionnaireRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Questionnaires{}, id)
	return result.Error
}
