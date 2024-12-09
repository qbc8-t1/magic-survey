package model

import (
	"errors"
	"strings"
	"time"
)

// validation errors
var (
	// create
	ErrInvalidQuestionIDCreate              = errors.New("questionID is required and must be greater than 0")
	ErrInvalidTitleCreate                   = errors.New("title is required and cannot be empty")
	ErrInvalidTypeCreate                    = errors.New("type is required and must be multioption or descriptive")
	ErrInvalidOrderCreate                   = errors.New("order is required and must be greater than 0")
	ErrInvalidFilePathCreate                = errors.New("filePath cannot be empty if provided")
	ErrInvalidDependsOnQuestionIDCreate     = errors.New("dependsOnQuestionID must be greater than 0 if provided")
	ErrInvalidDependsOnOptionIDCreate       = errors.New("dependsOnOptionID must be greater than 0 if provided")
	ErrDependsOnOptionWithoutQuestionCreate = errors.New("dependsOnQuestionID must be provided when DependsOnOptionID is provided")

	// update
	ErrAtLeatOneFieldNeededQuestion         = errors.New("at least one field must be provided for updating question")
	ErrInvalidQuestionIDUpdate              = errors.New("questionID must be greater than 0")
	ErrInvalidTitleUpdate                   = errors.New("title cannot be empty")
	ErrInvalidTypeUpdate                    = errors.New("type must be multioption or descriptive")
	ErrInvalidOrderUpdate                   = errors.New("order must be greater than 0")
	ErrInvalidFilePathUpdate                = errors.New("filePath cannot be empty")
	ErrInvalidDependsOnQuestionIDUpdate     = errors.New("dependsOnQuestionID must be greater than 0")
	ErrInvalidDependsOnOptionIDUpdate       = errors.New("dependsOnOptionID must be greater than 0")
	ErrDependsOnOptionWithoutQuestionUpdate = errors.New("dependsOnQuestionID must be provided when DependsOnOptionID is provided")
)

// QuestionsTypeEnum represents the questions_type_enum type in Postgres
type QuestionsTypeEnum string
type QuestionID uint

const (
	QuestionsTypeMultioption QuestionsTypeEnum = "multioption"
	QuestionsTypeDescriptive QuestionsTypeEnum = "descriptive"
)

// Question represents the questions table
type Question struct {
	ID                  QuestionID        `gorm:"primaryKey"`
	Title               string            `gorm:"size:255;not null"`
	Type                QuestionsTypeEnum `gorm:"not null"`
	QuestionnaireID     QuestionnaireID   `gorm:"not null"`
	Order               int               `gorm:"not null"`
	FilePath            *string           `gorm:"size:255;default:null"`
	DependsOnQuestionID *QuestionID       `gorm:"default:null"`
	DependsOnOptionID   *OptionID         `gorm:"default:null"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Questionnaire       Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	Options             *[]Option     `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
}

// CreateQuestionDTO represents the data needed to create a new question
type CreateQuestionDTO struct {
	Title               string            `json:"title" validate:"required"`
	Type                QuestionsTypeEnum `json:"type" validate:"required,oneof=multioption descriptive"`
	QuestionnaireID     QuestionnaireID   `json:"questionnaire_id" validate:"required"`
	Order               int               `json:"order" validate:"required"`
	FilePath            *string           `json:"file_path,omitempty" validate:"omitempty"`
	DependsOnQuestionID *QuestionID       `json:"depends_on_question_id,omitempty" validate:"omitempty"`
	DependsOnOptionID   *OptionID         `json:"depends_on_option_id,omitempty" validate:"omitempty"`
}

// UpdateQuestionDTO represents the data needed to update an existing question
type UpdateQuestionDTO struct {
	Title               *string            `json:"title,omitempty"`
	Type                *QuestionsTypeEnum `json:"type,omitempty"`
	QuestionnaireID     *QuestionnaireID   `json:"questionnaire_id,omitempty"`
	Order               *int               `json:"order,omitempty"`
	FilePath            *string            `json:"file_path,omitempty"`
	DependsOnQuestionID *QuestionID        `json:"depends_on_question_id,omitempty"`
	DependsOnOptionID   *OptionID          `json:"depends_on_option_id,omitempty"`
}

// QuestionResponse represents the question data returned in API responses
type QuestionResponse struct {
	ID                  QuestionID        `json:"id"`
	Title               string            `json:"title"`
	Type                QuestionsTypeEnum `json:"type"`
	QuestionnaireID     QuestionnaireID   `json:"questionnaire_id"`
	Order               int               `json:"order"`
	FilePath            *string           `json:"file_path"`
	DependsOnQuestionID *QuestionID       `json:"depends_on_question_id"`
	DependsOnOptionID   *OptionID         `json:"depends_on_option_id"`
}

// ToQuestionResponse maps a Question model to a QuestionResponseDTO
func ToQuestionResponse(question *Question) *QuestionResponse {
	return &QuestionResponse{
		ID:                  question.ID,
		Title:               question.Title,
		Type:                question.Type,
		QuestionnaireID:     question.QuestionnaireID,
		Order:               question.Order,
		FilePath:            question.FilePath,
		DependsOnQuestionID: question.DependsOnQuestionID,
		DependsOnOptionID:   question.DependsOnOptionID,
	}
}

func ToQuestionResponses(questions *[]Question) *[]QuestionResponse {
	questionResponses := make([]QuestionResponse, 0)
	for _, question := range *questions {
		questionResponses = append(questionResponses, QuestionResponse{
			ID:                  question.ID,
			Title:               question.Title,
			Type:                question.Type,
			QuestionnaireID:     question.QuestionnaireID,
			Order:               question.Order,
			FilePath:            question.FilePath,
			DependsOnQuestionID: question.DependsOnQuestionID,
			DependsOnOptionID:   question.DependsOnOptionID,
		})
	}

	return &questionResponses
}

// ToQuestionModel maps a CreateQuestionDTO to a Question model
func ToQuestionModel(questionDTO *CreateQuestionDTO) *Question {
	return &Question{
		Title:               questionDTO.Title,
		Type:                questionDTO.Type,
		QuestionnaireID:     questionDTO.QuestionnaireID,
		Order:               questionDTO.Order,
		FilePath:            questionDTO.FilePath,
		DependsOnQuestionID: questionDTO.DependsOnQuestionID,
		DependsOnOptionID:   questionDTO.DependsOnOptionID,
	}
}

// UpdateQuestionModel updates the fields of a Question model from an UpdateQuestionDTO
func UpdateQuestionModel(question *Question, questionDTO *UpdateQuestionDTO) {
	if questionDTO.Title != nil {
		question.Title = *questionDTO.Title
	}
	if questionDTO.Type != nil {
		question.Type = *questionDTO.Type
	}
	if questionDTO.QuestionnaireID != nil {
		question.QuestionnaireID = *questionDTO.QuestionnaireID
	}
	if questionDTO.Order != nil {
		question.Order = *questionDTO.Order
	}
	if questionDTO.FilePath != nil {
		question.FilePath = questionDTO.FilePath
	}
	if questionDTO.DependsOnQuestionID != nil {
		question.DependsOnQuestionID = questionDTO.DependsOnQuestionID
	}
	if questionDTO.DependsOnOptionID != nil {
		question.DependsOnOptionID = questionDTO.DependsOnOptionID
	}
}

// Validate validates a Question object
func (dto *CreateQuestionDTO) Validate() error {
	// Validate required fields
	if strings.TrimSpace(dto.Title) == "" {
		return ErrInvalidTitleCreate
	}
	// Validate Type
	if dto.Type != QuestionsTypeMultioption && dto.Type != QuestionsTypeDescriptive {
		return ErrInvalidTypeCreate
	}
	if dto.QuestionnaireID == 0 {
		return ErrInvalidQuestionnaireIDCreate
	}
	if dto.Order <= 0 {
		return ErrInvalidOrderCreate
	}
	// Validate FilePath (optional but cannot be empty if provided)
	if dto.FilePath != nil && strings.TrimSpace(*dto.FilePath) == "" {
		return ErrInvalidFilePathCreate
	}
	// Validate DependsOnQuestionID (optional but must be > 0 if provided)
	if dto.DependsOnQuestionID != nil && *dto.DependsOnQuestionID == 0 {
		return ErrInvalidDependsOnQuestionIDCreate
	}
	// Validate DependsOnOptionID (optional but must be > 0 if provided)
	if dto.DependsOnOptionID != nil && *dto.DependsOnOptionID == 0 {
		return ErrInvalidDependsOnOptionIDCreate
	}
	// If DependsOnOptionID is provided, DependsOnQuestionID must also be provided
	if dto.DependsOnOptionID != nil && dto.DependsOnQuestionID == nil {
		return ErrDependsOnOptionWithoutQuestionCreate
	}
	// Removed options validation
	return nil
}

func (dto *UpdateQuestionDTO) Validate() error {
	if dto.Title == nil && dto.Type == nil && dto.QuestionnaireID == nil && dto.Order == nil && dto.FilePath == nil && dto.DependsOnQuestionID == nil && dto.DependsOnOptionID == nil {
		return ErrAtLeatOneFieldNeededQuestion
	}
	// Validate Title (if provided)
	if dto.Title != nil && strings.TrimSpace(*dto.Title) == "" {
		return ErrInvalidTitleUpdate
	}
	// Validate Type (if provided)
	if dto.Type != nil {
		if *dto.Type != QuestionsTypeMultioption && *dto.Type != QuestionsTypeDescriptive {
			return ErrInvalidTypeUpdate
		}
	}
	// Validate QuestionnaireID (if provided)
	if dto.QuestionnaireID != nil && *dto.QuestionnaireID == 0 {
		return ErrInvalidQuestionnaireIDUpdate
	}
	// Validate Order (if provided)
	if dto.Order != nil && *dto.Order <= 0 {
		return ErrInvalidOrderUpdate
	}
	// Validate FilePath (if provided)
	if dto.FilePath != nil && strings.TrimSpace(*dto.FilePath) == "" {
		return ErrInvalidFilePathUpdate
	}
	// Validate DependsOnQuestionID (if provided)
	if dto.DependsOnQuestionID != nil && *dto.DependsOnQuestionID == 0 {
		return ErrInvalidDependsOnQuestionIDUpdate
	}
	// Validate DependsOnOptionID (if provided)
	if dto.DependsOnOptionID != nil && *dto.DependsOnOptionID == 0 {
		return ErrInvalidDependsOnOptionIDUpdate
	}
	// If DependsOnOptionID is provided, DependsOnQuestionID must also be provided
	if dto.DependsOnOptionID != nil && dto.DependsOnQuestionID == nil {
		return ErrDependsOnOptionWithoutQuestionUpdate
	}
	// Removed options validation
	return nil
}
