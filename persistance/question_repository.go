package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	// dependency injection
	db *gorm.DB
}

// NewQuestionRepository creates a new instance of QuestionRepository
func NewQuestionRepository(db *gorm.DB) domain_repository.IQuestionRepository {
	return &QuestionRepository{db: db}
}

// CreateQuestion adds a new question to the database
func (r *QuestionRepository) CreateQuestion(question *model.Question) error {
	return r.db.Create(&question).Error
}

// GetQuestionByID gets a question from database based on its ID
func (r *QuestionRepository) GetQuestionByID(id model.QuestionID) (*model.Question, error) {
	var question model.Question
	result := r.db.First(&question, id)
	return &question, result.Error
}

func (r *QuestionRepository) GetQuestionsByQuestionnaireID(questionnaireID model.QuestionnaireID) (*[]model.Question, error) {
	var questions []model.Question
	// Filter questions by QuestionnaireID
	result := r.db.Where("questionnaire_id = ?", questionnaireID).Find(&questions)
	return &questions, result.Error
}

// UpdateQuestion gets a question and updates it in database
func (r *QuestionRepository) UpdateQuestion(question *model.Question) error {
	return r.db.Save(&question).Error
}

// DeleteQuestion deletes a question from database
func (r *QuestionRepository) DeleteQuestion(id model.QuestionID) error {
	return r.db.Delete(&model.Question{}, id).Error
}

func (r *QuestionRepository) FindQuestionByQuestionIDAndQuestionnaireID(questionID uint, questionnaireID uint) (model.Question, error) {
	question := model.Question{}
	err := r.db.Find(&question, "id = ? and questionnaire_id = ?", questionID, questionnaireID).Error
	return question, err
}
