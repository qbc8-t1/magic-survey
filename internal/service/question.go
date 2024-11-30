package service

import domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"

type IQuestionService interface {
}

type QuestionService struct {
	// dependency injection
	repo domain_repository.IQuestionRepository
}

// NewQuestionService creates a new QuestionService object
func NewQuestionService(domain_repository.IQuestionRepository) IQuestionService {
	return &QuestionService{}
}
