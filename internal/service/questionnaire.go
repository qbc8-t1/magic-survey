package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

var (
	ErrQuestionnaireRetrieveFailed = errors.New("failed to retrieve questionnaire")
)

type IQuestionnaireService interface {
	GetQuestionnaireByID(questionnaireID uint) (model.Questionnaire, error)
}

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
