package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

// ISubmissionRepository interface defines the repository methods
type ISubmissionRepository interface {
	GetSubmissionByID(model.SubmissionID) (*model.Submission, error)
	CreateSubmission(submission *model.Submission) error
	UpdateSubmission(submission *model.Submission) error

	GetActiveSubmissionByUserID(userID model.UserID) (*model.Submission, error)
}
