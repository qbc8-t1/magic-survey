package model

import "time"

// SubmissionsStatusEnum represents the submissions_status_enum type in Postgres
type SubmissionsStatusEnum string

const (
	SubmissionsStatusAnswering SubmissionsStatusEnum = "answering"
	SubmissionsStatusSubmitted SubmissionsStatusEnum = "submitted"
	SubmissionsStatusCancelled SubmissionsStatusEnum = "cancelled"
	SubmissionsStatusClosed    SubmissionsStatusEnum = "closed"
)

// Submission represents the submissions table
type Submission struct {
	ID                     uint `gorm:"primaryKey"`
	QuestionnaireID        uint
	UserID                 uint
	Status                 SubmissionsStatusEnum `gorm:"type:submissions_status_enum;default:'answering'"`
	LastAnsweredQuestionID *uint
	SubmittedAt            *time.Time
	SpentMinutes           *int
	Questionnaire          Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	User                   User          `gorm:"foreignKey:UserID"`
	Answers                []Answer      `gorm:"foreignKey:SubmissionID"`
}
