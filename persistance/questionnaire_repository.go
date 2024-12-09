package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type QuestionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) domain_repository.IQuestionnaireRepository {
	return &QuestionnaireRepository{db: db}
}

func (r *QuestionnaireRepository) GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (*model.Questionnaire, error) {
	var questionnaire model.Questionnaire
	result := r.db.Preload("Questions.Options").Preload("Submissions.Answers.Option").First(&questionnaire, questionnaireID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &questionnaire, nil
}

func (r *QuestionnaireRepository) GetFirstQuestion(questionnaireID model.QuestionnaireID) (*model.Question, error) {
	var question model.Question
	result := r.db.Where("questionnaire_id = ?", questionnaireID).Order("\"order\" ASC").Preload("Options").First(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (r *QuestionnaireRepository) GetNextQuestion(questionnaireID model.QuestionnaireID, currentOrder int) (*model.Question, error) {
	var question model.Question
	result := r.db.Where("questionnaire_id = ? AND \"order\" > ?", questionnaireID, currentOrder).
		Order("\"order\" ASC").
		Preload("Options").
		First(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (r *QuestionnaireRepository) GetPreviousQuestion(questionnaireID model.QuestionnaireID, currentOrder int) (*model.Question, error) {
	var question model.Question
	result := r.db.Where("questionnaire_id = ? AND \"order\" < ?", questionnaireID, currentOrder).
		Order("\"order\" DESC").
		Preload("Options").
		First(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}
