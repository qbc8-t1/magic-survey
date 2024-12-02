package model

import (
	"errors"
	"strings"
	"time"
)

// Custom error messages
var (
	ErrInvalidUserID       = errors.New("UserID is required and must be greater than 0")
	ErrInvalidSubmissionID = errors.New("SubmissionID is required and must be greater than 0")
	ErrInvalidOptionID     = errors.New("OptionID must be greater than 0 if provided")
	ErrInvalidQuestionID   = errors.New("QuestionID is required and must be greater than 0")
	ErrInvalidAnswerText   = errors.New("AnswerText cannot be empty if provided")
	ErrConflictingFields   = errors.New("Either OptionID or AnswerText must be provided, but not both")
)

// Answer represents the answers table
type Answer struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"not null"`
	SubmissionID uint    `gorm:"not null"`
	QuestionID   uint    `gorm:"not null"`
	OptionID     *uint   `gorm:"default:null"`
	AnswerText   *string `gorm:"default:null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
	Question     Question   `gorm:"foreignKey:QuestionID"`
	Option       *Option    `gorm:"foreignKey:OptionID;default:null"`
}

// CreateANswerDTO represents the data needed to create a new answer
// TODO: check existence of user, submission, question and option in service
type CreateAnswerDTO struct {
	UserID       uint    `json:"user_id"`
	SubmissionID uint    `json:"submission_id"`
	QuestionID   uint    `json:"question_id"`
	OptionID     *uint   `json:"option_id,omitempty"`
	AnswerText   *string `json:"answer_text,omitempty"`
}

// UpdateAnswerDTO represents the data needed to update an existing answer
// TODO: check existence of option in service
type UpdateAnswerDTO struct {
	OptionID   *uint   `json:"option_id,omitempty"`
	AnswerText *string `json:"answer_text,omitempty"`
}

// AnswerResponse represents the answer data returned in API responses
type AnswerResponse struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	SubmissionID uint       `json:"submission_id"`
	QuestionID   uint       `json:"question_id"`
	OptionID     *uint      `json:"option_id"`
	AnswerText   *string    `json:"answer_text"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Submission   Submission `json:"submission"`
	Question     Question   `json:"question"`
	Option       *Option    `json:"option"`
}

// ToAnswerResponse maps an Answer model to an AnswerResponseDTO
func ToAnswerResponse(answer *Answer) *AnswerResponse {
	return &AnswerResponse{
		ID:           answer.ID,
		UserID:       answer.UserID,
		SubmissionID: answer.SubmissionID,
		QuestionID:   answer.QuestionID,
		OptionID:     answer.OptionID,
		AnswerText:   answer.AnswerText,
		CreatedAt:    answer.CreatedAt,
		UpdatedAt:    answer.UpdatedAt,
		Submission:   answer.Submission,
		Question:     answer.Question,
		Option:       answer.Option,
	}
}

func ToAnswerResponses(answers *[]Answer) *[]AnswerResponse {
	answerResponses := make([]AnswerResponse, 0)
	for _, answer := range *answers {
		answerResponses = append(answerResponses, AnswerResponse{
			ID:           answer.ID,
			UserID:       answer.UserID,
			SubmissionID: answer.SubmissionID,
			QuestionID:   answer.QuestionID,
			OptionID:     answer.OptionID,
			AnswerText:   answer.AnswerText,
			CreatedAt:    answer.CreatedAt,
			UpdatedAt:    answer.UpdatedAt,
			Submission:   answer.Submission,
			Question:     answer.Question,
			Option:       answer.Option,
		})
	}

	return &answerResponses
}

// ToAnswerModel maps a CreateAnswerDTO to a Answer model
func ToAnswerModel(answerDTO *CreateAnswerDTO) *Answer {
	return &Answer{
		UserID:       answerDTO.UserID,
		SubmissionID: answerDTO.SubmissionID,
		QuestionID:   answerDTO.QuestionID,
		OptionID:     answerDTO.OptionID,
		AnswerText:   answerDTO.AnswerText,
	}
}

func UpdateAnswerModel(answer *Answer, answerDTO *UpdateAnswerDTO) {
	if answerDTO.OptionID != nil {
		answer.OptionID = answerDTO.OptionID
	}
	if answerDTO.AnswerText != nil {
		answer.AnswerText = answerDTO.AnswerText
	}
}

// ValidateCreateAnswerDTO validates the fields in CreateAnswerDTO
func (dto *CreateAnswerDTO) Validate() error {
	// Validate required fields
	if dto.UserID == 0 {
		return ErrInvalidUserID
	}
	if dto.SubmissionID == 0 {
		return ErrInvalidSubmissionID
	}
	if dto.QuestionID == 0 {
		return ErrInvalidQuestionID
	}

	// Validate that either OptionID or AnswerText is provided, but not both
	if (dto.OptionID != nil && *dto.OptionID != 0) && (dto.AnswerText != nil && strings.TrimSpace(*dto.AnswerText) != "") {
		return ErrConflictingFields
	}

	// Validate answer text (optional but should be non-empty if provided)
	if dto.AnswerText != nil && strings.TrimSpace(*dto.AnswerText) == "" {
		return ErrInvalidAnswerText
	}

	// If OptionID is provided, ensure it's valid
	if dto.OptionID != nil && *dto.OptionID == 0 {
		return ErrInvalidOptionID
	}

	return nil
}

// ValidateUpdateAnswerDTO validates the fields in UpdateAnswerDTO
func (dto *UpdateAnswerDTO) Validate() error {
	// Validate that either OptionID or AnswerText is provided, but not both
	if (dto.OptionID != nil && *dto.OptionID != 0) && (dto.AnswerText != nil && strings.TrimSpace(*dto.AnswerText) != "") {
		return ErrConflictingFields
	}

	// Validate answer text (optional but should be non-empty if provided)
	if dto.AnswerText != nil && strings.TrimSpace(*dto.AnswerText) == "" {
		return ErrInvalidAnswerText
	}

	// If OptionID is provided, ensure it's valid
	if dto.OptionID != nil && *dto.OptionID == 0 {
		return ErrInvalidOptionID
	}

	return nil
}
