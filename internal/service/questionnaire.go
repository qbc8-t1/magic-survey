package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

var (
	ErrQuestionnaireNotFound                = errors.New("questionnaire not found")
	ErrQuestionnaireRetrieveFailed          = errors.New("failed to retrieve questionnaire")
	ErrNoNextQuestionAvailable              = errors.New("no next question available")
	ErrNoQuestionsInQuestionnaire           = errors.New("no questions available in this questionnaire")
	ErrQuestionDoesNotBelongToQuestionnaire = errors.New("question does not belong to the current questionnaire")
)

type IQuestionnaireService interface {
	GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (*model.Questionnaire, error)
}

type QuestionnaireService struct {
	repo domain_repository.IQuestionnaireRepository
}

func NewQuestionnaireService(repo domain_repository.IQuestionnaireRepository) *QuestionnaireService {
	return &QuestionnaireService{
		repo: repo,
	}
}

func (s *QuestionnaireService) GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (*model.Questionnaire, error) {
	return s.repo.GetQuestionnaireByID(questionnaireID)
}
