package service

import "github.com/QBC8-Team1/magic-survey/internal/domain/repository"

type IQuestionnaireService interface {
}

type QuestionnaireService struct {
	// dependency injection
	repo repository.QuestionnaireRepository
}

func NewQuestionnaireService(repository.IQuestionnaireRepository) IQuestionnaireService {
	return &QuestionnaireService{}
}
