package service

import (
	"errors"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

var (
	// General errors
	ErrOptionNotFound  = errors.New("option not found")
	ErrInvalidOptionID = errors.New("invalid option ID")

	// Creation errors
	ErrOptionCreateFailed = errors.New("failed to create option")

	// Caption is repetitive
	ErrOptionRepetitiveCaption = errors.New("option's caption is repetitive")

	// Update errors
	ErrOptionUpdateFailed = errors.New("failed to update option")

	// Delete errors
	ErrOptionDeleteFailed = errors.New("failed to delete option")

	// Retrieval errors
	ErrOptionRetrieveFailed = errors.New("failed to retrieve option")
)

type IOptionService interface {
	CreateOption(optionDTO *model.CreateOptionDTO) error
	GetOptionByID(id model.OptionID) (*model.OptionResponse, error)
	GetOptionsByQuestionID(quesionID model.QuestionID) (*[]model.OptionResponse, error)
	UpdateOption(id model.OptionID, optionDTO *model.UpdateOptionDTO) error
	DeleteOption(id model.OptionID) error
}

type OptionService struct {
	// dependency injection
	optionRepo   domain_repository.IOptionRepository
	questionRepo domain_repository.IQuestionRepository
}

// NewOptionService creates a new OptionService object
func NewOptionService(optionRepo domain_repository.IOptionRepository, questionRepo domain_repository.IQuestionRepository) *OptionService {
	return &OptionService{
		optionRepo:   optionRepo,
		questionRepo: questionRepo,
	}
}

func (s *OptionService) CreateOption(optionDTO *model.CreateOptionDTO) error {
	// Check if the Question exists
	questionID := optionDTO.QuestionID

	question, err := s.questionRepo.GetQuestionByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return ErrQuestionRetrieveFailed
	}

	// Convert DTO to model
	option := model.ToOptionModel(optionDTO)

	// Check whether caption is repetitive or not
	for _, opt := range *question.Options {
		if opt.Caption == option.Caption {
			return ErrOptionRepetitiveCaption
		}
	}

	option.CreatedAt = time.Now()

	// Create the option
	err = s.optionRepo.CreateOption(option)
	if err != nil {
		return ErrOptionCreateFailed
	}

	return nil
}

func (s *OptionService) UpdateOption(id model.OptionID, optionDTO *model.UpdateOptionDTO) error {
	// Check if the option exists
	existingOption, err := s.optionRepo.GetOptionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOptionNotFound
		}
		return ErrOptionRetrieveFailed
	}

	// If QuestionID is being updated, check if the new Question exists
	if optionDTO.QuestionID != nil {
		_, err := s.questionRepo.GetQuestionByID(*optionDTO.QuestionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrQuestionNotFound
			}
			return ErrQuestionRetrieveFailed
		}
	}

	existingOption.UpdatedAt = time.Now()
	model.UpdateOptionModel(existingOption, optionDTO)

	// Persist the changes
	err = s.optionRepo.UpdateOption(existingOption)
	if err != nil {
		return ErrOptionUpdateFailed
	}

	return nil
}

func (s *OptionService) GetOptionByID(id model.OptionID) (*model.OptionResponse, error) {
	option, err := s.optionRepo.GetOptionByID(id)
	if err != nil {
		return nil, ErrOptionNotFound
	}

	return model.ToOptionResponse(option), nil
}

func (s *OptionService) GetOptionsByQuestionID(questionID model.QuestionID) (*[]model.OptionResponse, error) {
	options, err := s.optionRepo.GetOptionsByQuestionID(questionID)
	if err != nil {
		return nil, ErrOptionNotFound
	}

	return model.ToOptionResponses(options), nil
}

func (s *OptionService) DeleteOption(id model.OptionID) error {
	// Check if the option exists
	_, err := s.optionRepo.GetOptionByID(id)
	if err != nil {
		return ErrOptionNotFound
	}

	// Delete the option
	err = s.optionRepo.DeleteOption(id)
	if err != nil {
		return ErrOptionDeleteFailed
	}

	return nil
}
