package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type SubmissionRepository struct {
	// dependency injection
	db *gorm.DB
}

// NewSubmissionRepository creates a new instance of SubmissionRepository
func NewSubmissionRepository(db *gorm.DB) domain_repository.ISubmissionRepository {
	return &SubmissionRepository{db: db}
}

// GetAnswerByID gets an answer from database based on its ID
func (r *SubmissionRepository) GetSubmissionByID(id model.SubmissionID) (*model.Submission, error) {
	var submission model.Submission
	result := r.db.First(&submission, id)
	return &submission, result.Error
}
