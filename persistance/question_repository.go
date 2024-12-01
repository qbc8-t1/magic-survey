package repository

import (
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	// dependency injection
	db *gorm.DB
}

// TODO: do times and uuid in service
// NewQuestionRepository creates a new instance of QuestionRepository
func NewQuestionRpository(db *gorm.DB) domain_repository.IQuestionRepository {
	return &QuestionRepository{db: db}
}

// CreateQuestion adds a new question to the database
func (r *QuestionRepository) CreateQuestion(question *model.Question) error {
	question.CreatedAt = time.Now()
	return r.db.Create(&question).Error
}

// GetQuestionByID gets a question from database based on its ID
func (r *QuestionRepository) GetQuestionByID(id uuid.UUID) (*model.Question, error) {
	var question model.Question
	result := r.db.First(&question, id)
	return &question, result.Error
}

func (r *QuestionRepository) GetQuestionsByID(ids []uuid.UUID) (*[]model.Question, error) {
	var questions []model.Question
	result := r.db.Find(&questions, ids)
	return &questions, result.Error
}

// UpdateQuestion gets a question and updates it in database
func (r *QuestionRepository) UpdateQuestion(question *model.Question) error {
	question.UpdatedAt = time.Now()
	return r.db.Save(&question).Error
}

// DeleteQuestion deletes a question from database
func (r *QuestionRepository) DeleteQuestion(id uuid.UUID) error {
	return r.db.Delete(&model.Question{}, id).Error
}
