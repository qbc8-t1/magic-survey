package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type OptionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) domain_repository.IOptionRepository {
	return &OptionRepository{db: db}
}

func (r *OptionRepository) CreateOption(option *model.Option) error {
	return r.db.Create(&option).Error
}

func (r *OptionRepository) GetOptionByID(id model.OptionID) (*model.Option, error) {
	var option model.Option
	result := r.db.First(&option, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &option, nil
}

func (r *OptionRepository) GetOptionsByQuestionID(questionID model.QuestionID) (*[]model.Option, error) {
	var options []model.Option
	result := r.db.Where("question_id = ?", questionID).Find(&options)
	if result.Error != nil {
		return nil, result.Error
	}
	return &options, nil
}

func (r *OptionRepository) UpdateOption(option *model.Option) error {
	return r.db.Save(&option).Error
}

func (r *OptionRepository) DeleteOption(id model.OptionID) error {
	return r.db.Delete(&model.Option{}, id).Error
}
