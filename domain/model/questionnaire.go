package model

import (
	"errors"
	"time"
)

var (
	// create
	ErrInvalidQuestionnaireIDCreate = errors.New("questionnaireID is required and must be greater than 0")

	// update
	ErrInvalidQuestionnaireIDUpdate = errors.New("questionnaireID must be greater than 0")
)

// QuestionnairesStatusEnum represents the questionnaires_status_enum type in Postgres
type QuestionnairesStatusEnum string
type QuestionnaireID uint

const (
	QuestionnaireStatusOpen      QuestionnairesStatusEnum = "open"
	QuestionnaireStatusClosed    QuestionnairesStatusEnum = "closed"
	QuestionnaireStatusCancelled QuestionnairesStatusEnum = "cancelled"
)

// QuestionnairesSequenceEnum represents the questionnaires_sequence_enum type in Postgres
type QuestionnairesSequenceEnum string

const (
	QuestionnaireTypeRandom     QuestionnairesSequenceEnum = "random"
	QuestionnaireTypeSequential QuestionnairesSequenceEnum = "sequential"
)

// QuestionnairesVisibilityEnum represents the questionnaires_visibility_enum type in Postgres
type QuestionnairesVisibilityEnum string

const (
	QuestionnaireVisibilityEverybody     QuestionnairesVisibilityEnum = "everybody"
	QuestionnaireVisibilityAdminAndOwner QuestionnairesVisibilityEnum = "admin_and_owner"
	QuestionnaireVisibilityNobody        QuestionnairesVisibilityEnum = "nobody"
)

// Questionnaire represents the questionnaires table
type Questionnaire struct {
	ID                         uint `gorm:"primaryKey"`
	OwnerID                    uint
	Status                     QuestionnairesStatusEnum `gorm:"type:questionnaires_status_enum;default:'open'"`
	CanSubmitFrom              *time.Time
	CanSubmitUntil             *time.Time
	MaxMinutesToResponse       int
	MaxMinutesToChangeAnswer   int
	MaxMinutesToGivebackAnswer int
	RandomOrSequential         QuestionnairesSequenceEnum `gorm:"type:questionnaires_sequence_enum"`
	CanBackToPreviousQuestion  bool
	Title                      string `gorm:"size:255"`
	MaxAllowedSubmissionsCount int
	AnswersVisibleFor          QuestionnairesVisibilityEnum `gorm:"type:questionnaires_visibility_enum"`
	CreatedAt                  time.Time
	Owner                      User         `gorm:"foreignKey:OwnerID"`
	Questions                  []Question   `gorm:"foreignKey:QuestionnaireID"`
	Submissions                []Submission `gorm:"foreignKey:QuestionnaireID"`
}
