package service

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

// CRUD

// createQuestionnaire()

// GetQuestionnaire()

// DeleteQuestionnaire()

// UpdateQuestionare()

type IQuestionnaireService interface {
	// repo domain_repository.IQuestionnaireRepository
	CreateQuestionnaire(questionnaire *model.Questionnaire) error
	GetQuestionnaireByID(id uint) (*model.Questionnaire, error)
	UpdateQuestionare(questionnaire *model.Questionnaire) error
	DeleteQuestionnaire(id uint) error
}

type QuestionnaireService struct {
	repo domain_repository.IQuestionnaireRepository
}

func NewQuestionnaireService(repo domain_repository.IQuestionnaireRepository) IQuestionnaireService {
	return &QuestionnaireService{repo: repo}
}

func (s *QuestionnaireService) CreateQuestionnaire(questionnaire *model.Questionnaire) error {
	return nil
}

func (s *QuestionnaireService) GetQuestionnaireByID(id uint) (*model.Questionnaire, error) {
	return nil, nil
}

func (s *QuestionnaireService) UpdateQuestionare(questionnaire *model.Questionnaire) error {
	return nil
}

func (s *QuestionnaireService) DeleteQuestionnaire(id uint) error {
	return nil
}
