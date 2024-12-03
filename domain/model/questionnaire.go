package model

import (
	"errors"
	"strings"
	"time"
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
	CanSubmitFrom              string `json:"can_submit_from,omitempty"`
	CanSubmitUntil             string `json:"can_submit_until,omitempty"`
	MaxMinutesToResponse       int    `json:"max_minutes_to_response,omitempty"`
	MaxMinutesToChangeAnswer   int    `json:"max_minutes_to_change_answer"`
	MaxMinutesToGivebackAnswer int    `json:"max_minutes_to_giveback_answer,omitempty"`

	RandomOrSequential         QuestionnairesSequenceEnum   `json:"random_or_sequential"`
	CanBackToPreviousQuestion  bool                         `json:"can_back_to_previous_question"`
	Title                      string                       `json:"title"`
	MaxAllowedSubmissionsCount int                          `json:"max_allowed_submissions_count"`
	AnswersVisibleFor          QuestionnairesVisibilityEnum `json:"answers_visible_for"`
}

type QuestionnaireResponse struct {
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
}

func ToQuestionnaireResponse(q *Questionnaire) *QuestionnaireResponse {
	return &QuestionnaireResponse{
		ID:                         q.ID,
		OwnerID:                    q.OwnerID,
		Status:                     q.Status,
		CanSubmitFrom:              q.CanSubmitFrom,
		CanSubmitUntil:             q.CanSubmitUntil,
		MaxMinutesToResponse:       q.MaxMinutesToResponse,
		MaxMinutesToChangeAnswer:   q.MaxMinutesToChangeAnswer,
		MaxMinutesToGivebackAnswer: q.MaxMinutesToGivebackAnswer,
		RandomOrSequential:         q.RandomOrSequential,
		CanBackToPreviousQuestion:  q.CanBackToPreviousQuestion,
		Title:                      q.Title,
		MaxAllowedSubmissionsCount: q.MaxAllowedSubmissionsCount,
		AnswersVisibleFor:          q.AnswersVisibleFor,
		CreatedAt:                  q.CreatedAt,
	}
}

func (dto CreateQuestionnaireDTO) ValidateAndMakeObject() (Questionnaire, error) {
	questionnaire := new(Questionnaire)

	// can_submit_from
	canSubmitFrom, err := time.Parse(time.DateTime, dto.CanSubmitFrom)
	if err != nil {
		return *questionnaire, errors.New("can_submit_from date is invalid. layout is this: 2006-01-02 15:04:05")
	}
	questionnaire.CanSubmitFrom = canSubmitFrom

	// can_submit_until
	canSubmitUntil, err := time.Parse(time.DateTime, dto.CanSubmitUntil)
	if err != nil {
		return *questionnaire, errors.New("can_submit_until date. layout is this: 2006-01-02 15:04:05")
	}
	questionnaire.CanSubmitUntil = canSubmitUntil

	// max_minutes_to_response
	if dto.MaxMinutesToResponse < 1 {
		return *questionnaire, errors.New("max_minutes_to_response is invalid")
	}
	questionnaire.MaxMinutesToResponse = dto.MaxMinutesToResponse

	// max_minutes_to_change_answer
	if dto.MaxMinutesToChangeAnswer < 1 {
		return *questionnaire, errors.New("max_minutes_to_change_answer is invalid")
	}
	questionnaire.MaxMinutesToChangeAnswer = dto.MaxMinutesToChangeAnswer

	// max_minutes_to_giveback_answer
	if dto.MaxMinutesToGivebackAnswer < 1 {
		return *questionnaire, errors.New("max_minutes_to_giveback_answer is invalid")
	}
	questionnaire.MaxMinutesToGivebackAnswer = dto.MaxMinutesToGivebackAnswer

	// random_or_sequential
	switch {
	case dto.RandomOrSequential == "rand" || dto.RandomOrSequential == "random":
		questionnaire.RandomOrSequential = "random"
	case dto.RandomOrSequential == "seq" || dto.RandomOrSequential == "sequential":
		questionnaire.RandomOrSequential = "sequential"
	default:
		return *questionnaire, errors.New("value of random_or_sequential field must be rand or seq")
	}

	// can_back_to_previous_question
	questionnaire.CanBackToPreviousQuestion = dto.CanBackToPreviousQuestion

	// title
	if len(strings.TrimSpace(dto.Title)) < 2 {
		return *questionnaire, errors.New("title is too short")
	}
	questionnaire.Title = dto.Title

	// max_allowed_submissions_count
	if dto.MaxAllowedSubmissionsCount < 1 {
		return *questionnaire, errors.New("max_allowed_submissions_count must be bigger than 0")
	}
	questionnaire.MaxAllowedSubmissionsCount = dto.MaxAllowedSubmissionsCount

	// answers_visible_for
	switch dto.AnswersVisibleFor {
	case "everybody":
	case "admin_and_owner":
	case "nobody":
	default:
		return *questionnaire, errors.New("value of answers_visible_for field must be one of these: everybody, admin_and_owner, nobody")
	}
	questionnaire.AnswersVisibleFor = dto.AnswersVisibleFor

	return *questionnaire, nil
}
