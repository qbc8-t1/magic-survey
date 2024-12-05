package service

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

type QuestionnaireService struct {
	repo domain_repository.IQuestionnaireRepository
}

func NewQuestionnaireService(repo domain_repository.IQuestionnaireRepository) *QuestionnaireService {
	return &QuestionnaireService{
		repo: repo,
	}
}

func (s *QuestionnaireService) GetQuestionnaireByID(questionnaireID uint) (model.Questionnaire, error) {
	return s.repo.GetQuestionnaireByID(questionnaireID)
}
