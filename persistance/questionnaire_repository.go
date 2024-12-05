package repository

import (
	"errors"

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

func (qr *QuestionnaireRepository) GetQuestionnaireByID(questionnnaireID uint) (model.Questionnaire, error) {
	questionnaire := new(model.Questionnaire)
	err := qr.db.First(questionnaire, "id = ?", questionnnaireID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Questionnaire{}, errors.New("questionnaire not found")
		}

		return model.Questionnaire{}, err
	}

	return *questionnaire, nil
}
