package model

import (
	"errors"
	"strings"
	"time"
)

// QuestionsTypeEnum represents the questions_type_enum type in Postgres
type QuestionsTypeEnum string

const (
	QuestionsTypeMultioption QuestionsTypeEnum = "multioption"
	QuestionsTypeDescriptive QuestionsTypeEnum = "descriptive"
)

// Question represents the questions table
type Question struct {
	ID                  uint              `gorm:"primaryKey"`
	Title               string            `gorm:"size:255;not null"`
	Type                QuestionsTypeEnum `gorm:"type:questions_type_enum;not null"`
	QuestionnaireID     uint
	Order               int
	FilePath            *string `gorm:"size:255;default:null"`
	DependsOnQuestionID *uint   `gorm:"default:null"`
	DependsOnOptionID   *uint   `gorm:"default:null"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Questionnaire       Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	Options             []Option      `gorm:"foreignKey:QuestionID"`
}

// CreateQuestionDTO represents the data needed to create a new question
type CreateQuestionDTO struct {
	Title               string            `json:"title" validate:"required"`
	Type                QuestionsTypeEnum `json:"type" validate:"required,oneof=multioption descriptive"`
	QuestionnaireID     uint              `json:"questionnaire_id" validate:"required"`
	Order               int               `json:"order" validate:"required"`
	FilePath            *string           `json:"file_path" validate:"omitempty"`
	DependsOnQuestionID *uint             `json:"depends_on_question_id" validate:"omitempty"`
	DependsOnOptionID   *uint             `json:"depends_on_option_id" validate:"omitempty"`
	Options             []Option          `json:"options,omitempty"`
}

// UpdateQuestionDTO represents the data needed to update an existing question
type UpdateQuestionDTO struct {
	Title               *string            `json:"title,omitempty"`
	Type                *QuestionsTypeEnum `json:"type,omitempty"`
	QuestionnaireID     *uint              `json:"questionnaire_id,omitempty"`
	Order               *int               `json:"order,omitempty"`
	FilePath            *string            `json:"file_path,omitempty"`
	DependsOnQuestionID *uint              `json:"depends_on_question_id,omitempty"`
	DependsOnOptionID   *uint              `json:"depends_on_option_id,omitempty"`
	Options             *[]Option          `json:"options,omitempty"`
}

// QuestionResponse represents the question data returned in API responses
type QuestionResponse struct {
	ID                  uint              `json:"id"`
	Title               string            `json:"title"`
	Type                QuestionsTypeEnum `json:"type"`
	QuestionnaireID     uint              `json:"questionnaire_id"`
	Order               int               `json:"order"`
	FilePath            *string           `json:"file_path"`
	DependsOnQuestionID *uint             `json:"depends_on_question_id"`
	DependsOnOptionID   *uint             `json:"depends_on_option_id"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	Options             []Option          `json:"options"`
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
		CreatedAt:           question.CreatedAt,
		UpdatedAt:           question.UpdatedAt,
		Options:             question.Options,
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
		Options:             questionDTO.Options,
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
	if questionDTO.Options != nil {
		question.Options = *questionDTO.Options
	}
}

// Validate validates a Question object
func (question *CreateQuestionDTO) Validate() error {
	// Validate Title
	if strings.TrimSpace(question.Title) == "" {
		return errors.New("title is required")
	}

	// Validate Type
	if question.Type != QuestionsTypeMultioption && question.Type != QuestionsTypeDescriptive {
		return errors.New("invalid question type")
	}

	// Validate QuestionnaireID
	if question.QuestionnaireID == 0 {
		return errors.New("questionnaire_id is required")
	}

	// Validate Order
	if question.Order <= 0 {
		return errors.New("order must be greater than zero")
	}

	// Conditional Validation for Options
	if question.Type == QuestionsTypeMultioption && len(question.Options) == 0 {
		return errors.New("options cannot be empty for multioption questions")
	}
	return nil
}

func (question *UpdateQuestionDTO) Validate() error {
	// Validate Title if provided
	if question.Title != nil && strings.TrimSpace(*question.Title) == "" {
		return errors.New("title cannot be empty")
	}

	// Validate Type if provided
	if question.Type != nil {
		if *question.Type != QuestionsTypeMultioption && *question.Type != QuestionsTypeDescriptive {
			return errors.New("invalid question type")
		}
	}

	// Validate QuestionnaireID if provided
	if question.QuestionnaireID != nil && *question.QuestionnaireID == 0 {
		return errors.New("questionnaire_id must be greater than zero")
	}

	// Validate Order if provided
	if question.Order != nil && *question.Order <= 0 {
		return errors.New("order must be greater than zero")
	}

	// Conditional Validation for Options if provided and Type is multioption
	if question.Type != nil && *question.Type == QuestionsTypeMultioption {
		if question.Options != nil && len(*question.Options) == 0 {
			return errors.New("options cannot be empty for multioption questions")
		}
	}

	return nil
}
