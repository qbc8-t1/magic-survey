package service

import "errors"

// service errors
var (
	// General errors
	ErrSubmissionNotFound      = errors.New("submission not found")
	ErrNoActiveSubmissionFound = errors.New("no active submission found")
	ErrInvalidSubmissionID     = errors.New("invalid submission ID")

	// Creation errors
	ErrSubmissionCreateFailed = errors.New("failed to create submission")

	// Retrieval errors
	ErrSubmissionRetrieveFailed = errors.New("failed to retrieve submission")

	// update errors
	ErrSubmissionUpdateFailed = errors.New("faild to update submission")
)
