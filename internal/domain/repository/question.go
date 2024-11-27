package repository

import "gorm.io/gorm"

type IQuestionRepository interface {
	// CreateQuestion() error
	// GetQuestion() error
	// UpdateQuestion() error
	// DeleteQuestion() error
}

type QuestionRepository struct {
	// dependency injection
	db *gorm.DB
}

func NewQuestionRpository(db *gorm.DB) IQuestionRepository {
	return &QuestionRepository{db: db}
}
