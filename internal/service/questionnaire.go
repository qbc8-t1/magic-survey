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

type QuestionnaireService struct {
	repo domain_repository.IQuestionnaireRepository
}

func NewQuestionnaireService(repo domain_repository.IQuestionnaireRepository) *QuestionnaireService {
	return &QuestionnaireService{repo: repo}
}

func CreateQuestionnaire(questionnaire *model.Questionnaire) error {
	return nil
}

func GetQuestionnaireByID(id uint) (*model.Questionnaire, error) {
	return nil, nil
}

func UpdateQuestionare(questionnaire *model.Questionnaire) error {
	return nil
}

func DeleteQuestionnaire(id uint) error {
	return nil
}
