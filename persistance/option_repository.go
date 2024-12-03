package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type OptionRepository struct {
	// dependency injection
	db *gorm.DB
}

// NewOptionRepository creates a new instance of OptionRepository
func NewOptionRepository(db *gorm.DB) domain_repository.IOptionRepository {
	return &OptionRepository{db: db}
}

// CreateOption adds a new option to the database
func (r *OptionRepository) CreateOption(option *model.Option) error {
	return r.db.Create(&option).Error
}

// GetOptionByID gets an option from database based on its ID
func (r *OptionRepository) GetOptionByID(id model.OptionID) (*model.Option, error) {
	var option model.Option
	result := r.db.First(&option, id)
	return &option, result.Error
}

func (r *OptionRepository) GetOptionsByQuestionID(questionID model.QuestionID) (*[]model.Option, error) {
	var options []model.Option
	// Filter options by QuestionID
	result := r.db.Where("question_id = ?", questionID).Find(&options)
	return &options, result.Error
}

// UpdateOption gets an option and updates it in database
func (r *OptionRepository) UpdateOption(option *model.Option) error {
	return r.db.Save(&option).Error
}

// DeleteOption deletes an option from database
func (r *OptionRepository) DeleteOption(id model.OptionID) error {
	return r.db.Delete(&model.Option{}, id).Error
}
