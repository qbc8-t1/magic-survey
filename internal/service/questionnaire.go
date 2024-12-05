package service

import "errors"

var (
	ErrQuestionnaireRetrieveFailed          = errors.New("failed to retrieve questionnaire")
	ErrNoNextQuestionAvailable              = errors.New("no next question available")
	ErrNoQuestionsInQuestionnaire           = errors.New("no questions available in this questionnaire")
	ErrQuestionDoesNotBelongToQuestionnaire = errors.New("question does not belong to the current questionnaire")
)
