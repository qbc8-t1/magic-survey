package model

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
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
	UserID                 UserID
	Status                 SubmissionsStatusEnum `gorm:"type:submissions_status_enum;default:'answering'"`
	CurrentQuestionID      *QuestionID           `gorm:"default:null"`
	LastAnsweredQuestionID *QuestionID           `gorm:"default:null"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	SubmittedAt            *time.Time    `gorm:"default:null"`
	SpentMinutes           *int          `gorm:"default:null"`
	QuestionOrder          []QuestionID  `gorm:"-" json:"-"`
	QuestionOrderRaw       []byte        `gorm:"type:jsonb;column:question_order"`
	Questionnaire          Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	User                   User          `gorm:"foreignKey:UserID"`
	Answers                []Answer      `gorm:"foreignKey:SubmissionID"`
}

// Hooks for marshaling/unmarshaling
func (s *Submission) BeforeSave(tx *gorm.DB) error {
	data, err := json.Marshal(s.QuestionOrder)
	if err != nil {
		return err
	}
	s.QuestionOrderRaw = data
	return nil
}

func (s *Submission) AfterFind(tx *gorm.DB) error {
	if len(s.QuestionOrderRaw) > 0 {
		var questionOrder []QuestionID
		if err := json.Unmarshal(s.QuestionOrderRaw, &questionOrder); err != nil {
			return err
		}
		s.QuestionOrder = questionOrder
	}
	return nil
}
