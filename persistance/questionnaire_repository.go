package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"

	"gorm.io/gorm"
)

type QuestionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) *QuestionnaireRepository {
	return &QuestionnaireRepository{db: db}
}

func (r *QuestionnaireRepository) GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (*model.Questionnaire, error) {
	var questionnaire model.Questionnaire
	result := r.db.First(&questionnaire, questionnaireID)
	return &questionnaire, result.Error
}

func (r *QuestionnaireRepository) GetFirstQuestion(questionnaireID model.QuestionnaireID) (*model.Question, error) {
	var question model.Question
	result := r.db.Where("questionnaire_id = ?", questionnaireID).Order("\"order\" ASC").First(&question)
	return &question, result.Error
}

func (r *QuestionnaireRepository) GetNextQuestion(questionnaireID model.QuestionnaireID, currentOrder int) (*model.Question, error) {
	var question model.Question
	result := r.db.Where("questionnaire_id = ? AND \"order\" > ?", questionnaireID, currentOrder).Order("\"order\" ASC").First(&question)
	return &question, result.Error
}

func (r *QuestionnaireRepository) GetPreviousQuestion(questionnaireID model.QuestionnaireID, currentOrder int) (*model.Question, error) {
	var question model.Question
	result := r.db.Where("questionnaire_id = ? AND \"order\" < ?", questionnaireID, currentOrder).Order("\"order\" DESC").First(&question)
	return &question, result.Error
}
