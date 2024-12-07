package service

import "errors"

var (
	ErrSubmissionRetrieveFailed = errors.New("failed to retrieve submission")
	ErrSubmissionCreateFailed   = errors.New("failed to create submission")
	ErrSubmissionUpdateFailed   = errors.New("faild to update submission")
	ErrNoActiveSubmissionFound  = errors.New("no active submission found")
)
