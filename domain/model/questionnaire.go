package model

import (
	"errors"
	"time"

	"github.com/QBC8-Team1/magic-survey/pkg/utils"
)

// QuestionnairesStatusEnum represents the questionnaires_status_enum type in Postgres
type QuestionnairesStatusEnum string

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
	CanSubmitFrom              time.Time
	CanSubmitUntil             time.Time
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

// CreateQuestionnaireDTO represents the data needed to create a new questionnaire
type CreateQuestionnaireDTO struct {
	// OwnerID              uint      `json:"owner_id" validate:"required"`
	CanSubmitFrom              time.Time `json:"can_submit_from,omitempty"`
	CanSubmitUntil             time.Time `json:"can_submit_until,omitempty"`
	MaxMinutesToResponse       int       `json:"max_minutes_to_response,omitempty"`
	MaxMinutesToChangeAnswer   int       `json:"max_minutes_to_change_answer"`
	MaxMinutesToGivebackAnswer int       `json:"max_minutes_to_give_back_answer,omitempty"`

	RandomOrSequential         QuestionnairesSequenceEnum   `json:"random_or_seq"`
	CanBackToPreviousQuestion  bool                         `json:"can_back_to_previous_question"`
	Title                      string                       `json:"title"`
	MaxAllowedSubmissionsCount int                          `json:"max_allowed_submission_count"`
	AnswersVisibleFor          QuestionnairesVisibilityEnum `json:"answers_visible_for"`
}

// UpdateUserDTO represents the data needed to update an existing user
type UpdateQuestionnaireDTO struct {
	FirstName    *string     `json:"first_name,omitempty"`
	LastName     *string     `json:"last_name,omitempty"`
	Email        *string     `json:"email,omitempty" validate:"email"`
	NationalCode *string     `json:"national_code,omitempty"`
	Password     *string     `json:"password,omitempty"`
	Gender       *GenderEnum `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
}

func ToQuestionnaireModel(dto *CreateQuestionnaireDTO) *Questionnaire {
	return &Questionnaire{
		// ID
		// OwnerID
		// Status
		CanSubmitFrom:              dto.CanSubmitFrom,
		CanSubmitUntil:             dto.CanSubmitUntil,
		MaxMinutesToResponse:       dto.MaxMinutesToResponse,
		MaxMinutesToChangeAnswer:   dto.MaxMinutesToChangeAnswer,
		MaxMinutesToGivebackAnswer: dto.MaxMinutesToGivebackAnswer,
		RandomOrSequential:         dto.RandomOrSequential,
		CanBackToPreviousQuestion:  dto.CanBackToPreviousQuestion,
		Title:                      dto.Title,
		MaxAllowedSubmissionsCount: dto.MaxAllowedSubmissionsCount,
		AnswersVisibleFor:          dto.AnswersVisibleFor,
		// CreatedAt                  time.Time
		// Owner                      User         `gorm:"foreignKey:OwnerID"`
		// Questions                  []Question   `gorm:"foreignKey:QuestionnaireID"`
		// Submissions                []Submission `gorm:"foreignKey:QuestionnaireID"`

	}
}

// Validate checks the User struct for common validation rules.
func (q *Questionnaire) Validate() error {
	// if strings.TrimSpace(u.FirstName) == "" {
	// 	return errors.New("first name is required")
	// }
	// if strings.TrimSpace(u.LastName) == "" {
	// 	return errors.New("last name is required")
	// }
	// if !utils.IsValidEmail(u.Email) {
	// 	return errors.New("invalid email format")
	// }
	// if len(u.NationalCode) != 10 || !utils.IsAllDigits(u.NationalCode) {
	// 	return errors.New("national code must be a 10-digit number")
	// }
	// if len(u.Password) < 6 {
	// 	return errors.New("password must be at least 6 characters long")
	// }

	// CanSubmitFrom:              dto.CanSubmitFrom,
	// CanSubmitUntil:             dto.CanSubmitUntil,
	// MaxMinutesToResponse:       dto.MaxMinutesToResponse,
	// MaxMinutesToChangeAnswer:   dto.MaxMinutesToChangeAnswer,
	// MaxMinutesToGivebackAnswer: dto.MaxMinutesToGivebackAnswer,
	// RandomOrSequential:         dto.RandomOrSequential,
	// CanBackToPreviousQuestion:  dto.CanBackToPreviousQuestion,
	// Title:                      dto.Title,
	// MaxAllowedSubmissionsCount: dto.MaxAllowedSubmissionsCount,
	// AnswersVisibleFor:          dto.AnswersVisibleFor,

	// ID                         uint `gorm:"primaryKey"`
	// OwnerID                    uint
	// Status                     QuestionnairesStatusEnum `gorm:"type:questionnaires_status_enum;default:'open'"`
	// CanSubmitFrom              time.Time
	// CanSubmitUntil             time.Time
	// MaxMinutesToResponse       int
	// MaxMinutesToChangeAnswer   int
	// MaxMinutesToGivebackAnswer int
	// RandomOrSequential         QuestionnairesSequenceEnum `gorm:"type:questionnaires_sequence_enum"`
	// CanBackToPreviousQuestion  bool
	// Title                      string `gorm:"size:255"`
	// MaxAllowedSubmissionsCount int
	// AnswersVisibleFor          QuestionnairesVisibilityEnum `gorm:"type:questionnaires_visibility_enum"`
	// CreatedAt                  time.Time
	// Owner                      User         `gorm:"foreignKey:OwnerID"`
	// Questions                  []Question   `gorm:"foreignKey:QuestionnaireID"`
	// Submissions                []Submission `gorm:"foreignKey:QuestionnaireID"`

	// utils.CanSubmitFrom
	if utils.IsValidDate(q.CanSubmitFrom) == "" {
		return errors.New("first name is required")
	}

	return nil
}
