package service

import "github.com/QBC8-Team1/magic-survey/internal/domain/repository"

type IQuestionService interface {
}

type QuestionService struct {
	// dependency injection
	repo repository.QuestionRepository
}

func NewQuestionService(repository.IQuestionRepository) IQuestionService {
	return &QuestionService{}
}
