package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	// dependency injection
	db *gorm.DB
}

// NewAnswerRepository creates a new instance of AnswerRepository
func NewAnswerRepository(db *gorm.DB) domain_repository.IAnswerRepository {
	return &AnswerRepository{db: db}
}

// CreateAnswer adds a new answer to the database
func (r *AnswerRepository) CreateAnswer(answer *model.Answer) error {
	return r.db.Create(&answer).Error
}

// GetAnswerByID gets an answer from database based on its ID
func (r *AnswerRepository) GetAnswerByID(id model.AnswerID) (*model.Answer, error) {
	var answer model.Answer
	result := r.db.First(&answer, id)
	return &answer, result.Error
}

// UpdateAnswer gets an answer and updates it in database
func (r *AnswerRepository) UpdateAnswer(answer *model.Answer) error {
	return r.db.Save(&answer).Error
}

// DeleteAnswer deletes an answer from database
func (r *AnswerRepository) DeleteAnswer(id model.AnswerID) error {
	return r.db.Delete(&model.Answer{}, id).Error
}
