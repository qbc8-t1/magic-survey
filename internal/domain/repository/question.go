package repository

import "gorm.io/gorm"

type IQuestionRepository interface {
	Add() error
	Get() error
	Update() error
	Delete() error
}

type QuestionRepository struct {
	// dependency injection
	db *gorm.DB
}

func NewQuestionRpository(db *gorm.DB) IQuestionRepository {
	return &QuestionRepository{db: db}
}
