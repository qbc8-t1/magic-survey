package repository

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

const PAGE_SIZE = 1

type QuestionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) domain_repository.IQuestionnaireRepository {
	return &QuestionnaireRepository{db: db}
}

func (r *QuestionnaireRepository) GetUserQuestionnairesCount(userID model.UserID) (int64, error) {
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

func (r *QuestionnaireRepository) UpdateQuestionaire(questionnnaireID model.QuestionnaireID, updateData *model.Questionnaire) error {
	return r.db.Model(&model.Questionnaire{}).Where("id = ?", questionnnaireID).Updates(updateData).Error
}

func (r *QuestionnaireRepository) DeleteQuestionnaire(questionnnaireID model.QuestionnaireID) error {
	return r.db.Delete(&model.Questionnaire{}, questionnnaireID).Error
}

func (qr *QuestionnaireRepository) GetQuestionnaireByID(questionnnaireID model.QuestionnaireID) (model.Questionnaire, error) {
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

func (qr *QuestionnaireRepository) GetQuestionnairesByOwnerID(ownerID model.UserID, page int) ([]domain_repository.Questionnaire, error) {
	var questionnaires []domain_repository.Questionnaire

	offset := (page - 1) * PAGE_SIZE

	err := qr.db.Model(&model.Questionnaire{}).
		Preload("Questions").
		Preload("Questions.Options").
		Where("owner_id = ?", ownerID).
		Limit(PAGE_SIZE).
		Offset(offset).
		Find(&questionnaires).Error

	return questionnaires, err
}
