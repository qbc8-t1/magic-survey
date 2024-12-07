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

func (r *QuestionnaireRepository) GetUserQuestionnairesCount(userID uint) (int64, error) {
	var count int64
	result := r.db.Model(&model.Questionnaire{}).Where("owner_id = ?", userID).Count(&count)
	return count, result.Error
}

func (r *QuestionnaireRepository) CreateQuestionnaire(questionnaire *model.Questionnaire) (model.Questionnaire, error) {
	err := r.db.Create(questionnaire).Error
	if err != nil {
		return model.Questionnaire{}, err
	}

	return *questionnaire, err
}

func (r *QuestionnaireRepository) GetQuestionnaireByID(id uint) (*model.Questionnaire, error) {
	return nil, nil
}

func (r *QuestionnaireRepository) UpdateQuestionaire(questionnaire *model.Questionnaire) error {
	return r.db.Save(questionnaire).Error
}

func (r *QuestionnaireRepository) DeleteQuestionnaire(id uint) error {
	return nil
}
