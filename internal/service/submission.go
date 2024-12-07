package service

import "errors"

var (
	ErrSubmissionNotFound  = errors.New("submission not found")
	ErrInvalidSubmissionID = errors.New("invalid submission ID")

	ErrSubmissionRetrieveFailed = errors.New("failed to retrieve submission")
	ErrSubmissionCreateFailed   = errors.New("failed to create submission")
	ErrSubmissionUpdateFailed   = errors.New("faild to update submission")
	ErrNoActiveSubmissionFound  = errors.New("no active submission found")
)
