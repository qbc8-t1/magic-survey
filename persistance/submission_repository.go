package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) domain_repository.ISubmissionRepository {
	return &SuperadminRepository{db: db}
}

func (s *SuperadminRepository) GetSubmissionByID(submissionID uint) (model.Submission, error) {
	submission := new(model.Submission)
	result := s.db.First(submission, "id = ?", submissionID)
	return *submission, result.Error
}
