package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidSubmissionIDCreate = errors.New("submissionID is required and must be greater than 0")
)

// SubmissionsStatusEnum represents the submissions_status_enum type in postgres
type SubmissionsStatusEnum string
type SubmissionID uint

const (
	SubmissionsStatusAnswering SubmissionsStatusEnum = "answering"
	SubmissionsStatusSubmitted SubmissionsStatusEnum = "submitted"
	SubmissionsStatusCancelled SubmissionsStatusEnum = "cancelled"
	SubmissionsStatusClosed    SubmissionsStatusEnum = "closed"
)

// Submission represents the submissions table
type Submission struct {
	ID                     SubmissionID `gorm:"primaryKey"`
	QuestionnaireID        QuestionnaireID
	UserID                 UserId
	Status                 SubmissionsStatusEnum `gorm:"type:submissions_status_enum;default:'answering'"`
	LastAnsweredQuestionID *QuestionID
	SubmittedAt            *time.Time
	SpentMinutes           *int
	Questionnaire          Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	User                   User          `gorm:"foreignKey:UserID"`
	Answers                []Answer      `gorm:"foreignKey:SubmissionID"`
}
