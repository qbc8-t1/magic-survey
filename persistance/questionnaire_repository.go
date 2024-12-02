package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type QuestionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) domain_repository.IQuestionnaireRepository {
	return &QuestionnaireRepository{db: db}
}

func (r *QuestionnaireRepository) CreateQuestionnaire(questionnaire *model.Questionnaire) error {
	return nil
}

func (r *QuestionnaireRepository) GetQuestionnaireByID(id uint) (*model.Questionnaire, error) {
	return nil, nil
}

func (r *QuestionnaireRepository) UpdateQuestionare(questionnaire *model.Questionnaire) error {
	return nil
}

func (r *QuestionnaireRepository) DeleteQuestionnaire(id uint) error {
	return nil
}
