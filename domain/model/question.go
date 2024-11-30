package model

import "time"

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
