package model

import (
	"time"
)

// Answer represents the answers table
type Answer struct {
	ID           uint `gorm:"primaryKey"`
	SubmissionID uint
	QuestionID   uint
	UserID       uint
	OptionID     *uint
	AnswerText   *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
	Question     Question   `gorm:"foreignKey:QuestionID"`
	Option       *Option    `gorm:"foreignKey:OptionID"`
}
