package model

import (
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
	Title               string            `gorm:"size:255"`
	Type                QuestionsTypeEnum `gorm:"type:questions_type_enum"`
	QuestionnaireID     uint
	Order               int
	FilePath            string `gorm:"size:255"`
	DependsOnQuestionID *uint
	DependsOnOptionID   *uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Questionnaire       Questionnaire `gorm:"foreignKey:QuestionnaireID"`
	Options             []Option      `gorm:"foreignKey:QuestionID"`
}

// CreateQuestionDTO represents the data needed to create a new question
type CreateQuestionDTO struct {
	Title               string            `json:"title" validate:"required"`
	Type                QuestionsTypeEnum `json:"type" validate:"required"`
	QuestionnaireID     uint              `json:"questionnaire_id" validate:"required"`
	Order               int               `json:"order" validate:"required"`
	FilePath            string            `json:"file_path" validate:"required"`
	DependsOnQuestionID *uint
	DependsOnOptionID   *uint
	Options             []Option
}

// UpdateQuestionDTO represents the data needed to update an existing question
type UpdateQuestionDTO struct {
	Title           string            `json:"title,omitempty"`
	Type            QuestionsTypeEnum `json:"type,omitempty"`
	QuestionnaireID uint              `json:"questionnaire_id,omitempty"`
	Order           int               `json:"order,omitempty"`
	FilePath        string            `json:"file_path,omitempty"`
}

// QuestionResponse represents the question data returned in API responses
type QuestionResponse struct {
	ID                  uint
	Title               string
	Type                QuestionsTypeEnum
	QuestionnaireID     uint
	Order               int
	FilePath            string
	DependsOnQuestionID *uint
	DependsOnOptionID   *uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Questionnaire       Questionnaire
	Options             []Option
}

// ToQuestionResponse maps a Question model to a QuestionResponseDTO
func ToQuestionResponse(question *Question) *QuestionResponse {
	return &QuestionResponse{
		ID:              question.ID,
		Title:           question.Title,
		Type:            question.Type,
		QuestionnaireID: question.QuestionnaireID,
	}
}

// ToQuestionModel maps a CreateQuestionDTO to a Question model
func ToQuestionModel(questionDTO *CreateQuestionDTO) *Question {
	return &Question{
		Title:           questionDTO.Title,
		Type:            questionDTO.Type,
		QuestionnaireID: questionDTO.QuestionnaireID,
		Order:           questionDTO.Order,
		FilePath:        questionDTO.FilePath,
	}
}

// UpdateQuestionModel updates the fields of a Qestion model from an UpdateQuestionDTO
// func UpdateQuestionModel(question *Question, questionDTO *UpdadteQuestionDTO) {

// }
