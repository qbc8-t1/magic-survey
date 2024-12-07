package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type ISubmissionRepository interface {
	GetSubmissionByID(submissionID uint) (model.Submission, error)
}
