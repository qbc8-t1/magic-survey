package repository

import (
	"fmt"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) domain_repository.IQuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) CreateQuestion(question *model.Question) error {
	return r.db.Create(&question).Error
}

func (r *QuestionRepository) GetQuestionByID(id model.QuestionID) (*model.Question, error) {
	var question model.Question
	result := r.db.Preload("Options").First(&question, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &question, nil
}

func (r *QuestionRepository) GetQuestionsByQuestionnaireID(questionnaireID model.QuestionnaireID) (*[]model.Question, error) {
	var questions []model.Question
	result := r.db.Where("questionnaire_id = ?", questionnaireID).Preload("Options").Find(&questions)
	if result.Error != nil {
		return nil, result.Error
	}
	return &questions, nil
}

func (r *QuestionRepository) UpdateQuestion(question *model.Question) error {
	return r.db.Save(&question).Error
}

func (r *QuestionRepository) DeleteQuestion(id model.QuestionID) error {
	return r.db.Delete(&model.Question{}, id).Error
}

func (repo *QuestionRepository) GetUnansweredQuestions(questionnaireID model.QuestionnaireID, submissionID model.SubmissionID) (*[]model.Question, error) {
	var unansweredQuestions []model.Question

	subquery := repo.db.
		Table("answers").
		Select("question_id").
		Where("submission_id = ?", submissionID)

	err := repo.db.Where("questionnaire_id = ? AND id NOT IN (?)", questionnaireID, subquery).
		Preload("Options").
		Find(&unansweredQuestions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve unanswered questions: %w", err)
	}

	return &unansweredQuestions, nil
}

func (r *QuestionRepository) GetQuestionByQuestionIDAndQuestionnaireID(questionID model.QuestionID, questionnaireID model.QuestionnaireID) (*model.Question, error) {
	var question model.Question
	err := r.db.Where("id = ? AND questionnaire_id = ?", questionID, questionnaireID).
		Preload("Options").
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}
