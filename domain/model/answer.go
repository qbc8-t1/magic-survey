package model

import (
	"errors"
	"strings"
	"time"
)

type AnswerID uint

// Custom error messages
var (
	ErrInvalidAnswerText          = errors.New("answerText cannot be empty if provided")
	ErrConflictingFields          = errors.New("either optionID or answerText must be provided, but not both")
	ErrAtLeatOneFieldNeededAnswer = errors.New("at least one field must be provided for updating answer")
)

// Answer represents the answers table
type Answer struct {
	ID           AnswerID `gorm:"primaryKey"`
	UserID       UserId
	SubmissionID uint
	QuestionID   QuestionID
	OptionID     *OptionID `gorm:"default:null"`
	AnswerText   *string   `gorm:"default:null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
	Question     Question   `gorm:"foreignKey:QuestionID"`
	Option       *Option    `gorm:"foreignKey:OptionID;default:null"`
}

// CreateANswerDTO represents the data needed to create a new answer
type CreateAnswerDTO struct {
	UserID       UserId     `json:"user_id"`
	SubmissionID uint       `json:"submission_id"`
	QuestionID   QuestionID `json:"question_id"`
	OptionID     *OptionID  `json:"option_id,omitempty"`
	AnswerText   *string    `json:"answer_text,omitempty"`
}

// UpdateAnswerDTO represents the data needed to update an existing answer
type UpdateAnswerDTO struct {
	OptionID   *OptionID `json:"option_id,omitempty"`
	AnswerText *string   `json:"answer_text,omitempty"`
}

// AnswerResponse represents the answer data returned in API responses
type AnswerResponse struct {
	ID           AnswerID   `json:"id"`
	UserID       UserId     `json:"user_id"`
	SubmissionID uint       `json:"submission_id"`
	QuestionID   QuestionID `json:"question_id"`
	OptionID     *OptionID  `json:"option_id"`
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
		ID:           AnswerID(answer.ID),
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

type AnswerSummaryResponse struct {
	SubmissionID uint      `json:"submission_id"`
	OptionID     *OptionID `json:"option_id"`
	AnswerText   *string   `json:"answer_text"`
	Option       *Option   `json:"option"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToAnswerSummaryResponses(answers *[]Answer) *[]AnswerSummaryResponse {
	answerResponses := make([]AnswerSummaryResponse, 0, len(*answers))
	for _, answer := range *answers {
		answerResponses = append(answerResponses, AnswerSummaryResponse{
			SubmissionID: answer.SubmissionID,
			OptionID:     answer.OptionID,
			AnswerText:   answer.AnswerText,
			Option:       answer.Option,
			CreatedAt:    answer.CreatedAt,
		})
	}

	return &answerResponses
}

func ToAnswerResponses(answers *[]Answer) *[]AnswerResponse {
	answerResponses := make([]AnswerResponse, 0)
	for _, answer := range *answers {
		answerResponses = append(answerResponses, AnswerResponse{
			ID:           AnswerID(answer.ID),
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
	if dto.OptionID == nil && dto.AnswerText == nil {
		return ErrAtLeatOneFieldNeededAnswer
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
