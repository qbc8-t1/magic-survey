package model

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidQuestionID              = errors.New("questionID is required and must be greater than 0")
	ErrInvalidTitle                   = errors.New("title is required and cannot be empty")
	ErrInvalidType                    = errors.New("type is required and must be 'multioption' or 'descriptive'")
	ErrInvalidOrder                   = errors.New("order is required and must be greater than 0")
	ErrInvalidFilePath                = errors.New("filePath cannot be empty if provided")
	ErrInvalidDependsOnQuestionID     = errors.New("dependsOnQuestionID must be greater than 0 if provided")
	ErrInvalidDependsOnOptionID       = errors.New("dependsOnOptionID must be greater than 0 if provided")
	ErrDependsOnOptionWithoutQuestion = errors.New("dependsOnQuestionID must be provided when DependsOnOptionID is provided")
	ErrAtLeatOneFieldNeededQuestion   = errors.New("at least one field must be provided for updating question")
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
	Type                QuestionsTypeEnum `gorm:"type:questions_type_enum;not null"`
	QuestionnaireID     QuestionnaireID   `gorm:"not null"`
	Order               int               `gorm:"not null"`
	FilePath            *string           `gorm:"size:255;default:null"`
	DependsOnQuestionID *uint             `gorm:"default:null"`
	DependsOnOptionID   *uint             `gorm:"default:null"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Questionnaire       Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	Options             *[]Option     `gorm:"foreignKey:QuestionID"`
}

// CreateQuestionDTO represents the data needed to create a new question
type CreateQuestionDTO struct {
	Title               string            `json:"title" validate:"required"`
	Type                QuestionsTypeEnum `json:"type" validate:"required,oneof=multioption descriptive"`
	QuestionnaireID     QuestionnaireID   `json:"questionnaire_id" validate:"required"`
	Order               int               `json:"order" validate:"required"`
	FilePath            *string           `json:"file_path,omitempty" validate:"omitempty"`
	DependsOnQuestionID *uint             `json:"depends_on_question_id,omitempty" validate:"omitempty"`
	DependsOnOptionID   *uint             `json:"depends_on_option_id,omitempty" validate:"omitempty"`
}

// UpdateQuestionDTO represents the data needed to update an existing question
type UpdateQuestionDTO struct {
	Title               *string            `json:"title,omitempty"`
	Type                *QuestionsTypeEnum `json:"type,omitempty"`
	QuestionnaireID     *QuestionnaireID   `json:"questionnaire_id,omitempty"`
	Order               *int               `json:"order,omitempty"`
	FilePath            *string            `json:"file_path,omitempty"`
	DependsOnQuestionID *uint              `json:"depends_on_question_id,omitempty"`
	DependsOnOptionID   *uint              `json:"depends_on_option_id,omitempty"`
}

// QuestionResponse represents the question data returned in API responses
type QuestionResponse struct {
	ID                  QuestionID        `json:"id"`
	Title               string            `json:"title"`
	Type                QuestionsTypeEnum `json:"type"`
	QuestionnaireID     QuestionnaireID   `json:"questionnaire_id"`
	Order               int               `json:"order"`
	FilePath            *string           `json:"file_path"`
	DependsOnQuestionID *uint             `json:"depends_on_question_id"`
	DependsOnOptionID   *uint             `json:"depends_on_option_id"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	Options             *[]Option         `json:"options"`
}

// ToQuestionResponse maps a Question model to a QuestionResponseDTO
func ToQuestionResponse(question *Question) *QuestionResponse {
	return &QuestionResponse{
		ID:                  QuestionID(question.ID),
		Title:               question.Title,
		Type:                question.Type,
		QuestionnaireID:     question.QuestionnaireID,
		Order:               question.Order,
		FilePath:            question.FilePath,
		DependsOnQuestionID: question.DependsOnQuestionID,
		DependsOnOptionID:   question.DependsOnOptionID,
		CreatedAt:           question.CreatedAt,
		UpdatedAt:           question.UpdatedAt,
		Options:             question.Options,
	}
}

func ToQuestionResponses(questions *[]Question) *[]QuestionResponse {
	questionResponses := make([]QuestionResponse, 0)
	for _, question := range *questions {
		questionResponses = append(questionResponses, QuestionResponse{
			ID:                  QuestionID(question.ID),
			Title:               question.Title,
			Type:                question.Type,
			QuestionnaireID:     question.QuestionnaireID,
			Order:               question.Order,
			FilePath:            question.FilePath,
			DependsOnQuestionID: question.DependsOnQuestionID,
			DependsOnOptionID:   question.DependsOnOptionID,
			CreatedAt:           question.CreatedAt,
			UpdatedAt:           question.UpdatedAt,
			Options:             question.Options,
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
		return ErrInvalidTitle
	}
	// Validate Type
	if dto.Type != QuestionsTypeMultioption && dto.Type != QuestionsTypeDescriptive {
		return ErrInvalidType
	}
	if dto.QuestionnaireID == 0 {
		return ErrInvalidQuestionnaireID
	}
	if dto.Order <= 0 {
		return ErrInvalidOrder
	}
	// Validate FilePath (optional but cannot be empty if provided)
	if dto.FilePath != nil && strings.TrimSpace(*dto.FilePath) == "" {
		return ErrInvalidFilePath
	}
	// Validate DependsOnQuestionID (optional but must be > 0 if provided)
	if dto.DependsOnQuestionID != nil && *dto.DependsOnQuestionID == 0 {
		return ErrInvalidDependsOnQuestionID
	}
	// Validate DependsOnOptionID (optional but must be > 0 if provided)
	if dto.DependsOnOptionID != nil && *dto.DependsOnOptionID == 0 {
		return ErrInvalidDependsOnOptionID
	}
	// If DependsOnOptionID is provided, DependsOnQuestionID must also be provided
	if dto.DependsOnOptionID != nil && dto.DependsOnQuestionID == nil {
		return ErrDependsOnOptionWithoutQuestion
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
		return ErrInvalidTitle
	}
	// Validate Type (if provided)
	if dto.Type != nil {
		if *dto.Type != QuestionsTypeMultioption && *dto.Type != QuestionsTypeDescriptive {
			return ErrInvalidType
		}
	}
	// Validate QuestionnaireID (if provided)
	if dto.QuestionnaireID != nil && *dto.QuestionnaireID == 0 {
		return ErrInvalidQuestionnaireID
	}
	// Validate Order (if provided)
	if dto.Order != nil && *dto.Order <= 0 {
		return ErrInvalidOrder
	}
	// Validate FilePath (if provided)
	if dto.FilePath != nil && strings.TrimSpace(*dto.FilePath) == "" {
		return ErrInvalidFilePath
	}
	// Validate DependsOnQuestionID (if provided)
	if dto.DependsOnQuestionID != nil && *dto.DependsOnQuestionID == 0 {
		return ErrInvalidDependsOnQuestionID
	}
	// Validate DependsOnOptionID (if provided)
	if dto.DependsOnOptionID != nil && *dto.DependsOnOptionID == 0 {
		return ErrInvalidDependsOnOptionID
	}
	// If DependsOnOptionID is provided, DependsOnQuestionID must also be provided
	if dto.DependsOnOptionID != nil && dto.DependsOnQuestionID == nil {
		return ErrDependsOnOptionWithoutQuestion
	}
	// Removed options validation
	return nil
}
