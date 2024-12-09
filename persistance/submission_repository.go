package repository

import (
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) domain_repository.ISubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) GetSubmissionByID(id model.SubmissionID) (*model.Submission, error) {
	var submission model.Submission
	result := r.db.Preload("Questionnaire.Questions.Options").
		Preload("Answers.Option").
		Preload("Answers.Question.Options").
		First(&submission, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &submission, nil
}

func (r *SubmissionRepository) CreateSubmission(submission *model.Submission) error {
	submission.CreatedAt = time.Now()
	return r.db.Create(submission).Error
}

func (r *SubmissionRepository) UpdateSubmission(submission *model.Submission) error {
	submission.UpdatedAt = time.Now()
	return r.db.Save(submission).Error
}

func (r *SubmissionRepository) GetActiveSubmissionByUserID(userID model.UserId) (*model.Submission, error) {
	var submission model.Submission
	result := r.db.Preload("Questionnaire.Questions.Options").
		Preload("Answers.Option").
		Preload("Answers.Question.Options").
		Where("user_id = ? AND status = ?", userID, model.SubmissionsStatusAnswering).
		First(&submission)
	if result.Error != nil {
		return nil, result.Error
	}
	return &submission, nil
}
