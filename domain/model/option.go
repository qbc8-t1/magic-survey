package model

import (
	"errors"
	"strings"
)

var (
	ErrInvalidOptionID            = errors.New("optionID must be greater than 0 if provided")
	ErrInvalidCaption             = errors.New("caption is required and cannot be empty or more than 255 characters")
	ErrInvalidIsCorrect           = errors.New("is_correct must be 'true' or 'false' if provided")
	ErrAtLeatOneFieldNeededOption = errors.New("at least one field must be provided for updating option")
)

type OptionID uint

// Option represents the options table
type Option struct {
	ID         OptionID   `gorm:"primaryKey"`
	QuestionID QuestionID `gorm:"not null"`
	Order      int        `gorm:"not null"`
	Caption    string     `gorm:"size:255;not null"`
	IsCorrect  *bool      `gorm:"default:false"`
	Question   Question   `gorm:"foreignKey:QuestionID"`
}

// CreateOptionDTO represents the data needed to create a new option
type CreateOptionDTO struct {
	QuestionID QuestionID `json:"question_id" validate:"required"`
	Order      int        `json:"order" validate:"required"`
	Caption    string     `json:"caption" validate:"required"`
	IsCorrect  *bool      `json:"is_correct,omitempty" validate:"omitempty"`
}

// UpdateOptionDTO represents the data needed to update an existing option
type UpdateOptionDTO struct {
	QuestionID *QuestionID `json:"question_id,omitempty"`
	Order      *int        `json:"order,omitempty"`
	Caption    *string     `json:"caption,omitempty"`
	IsCorrect  *bool       `json:"is_correct,omitempty"`
}

// OptionResponse represents the option data returned in API responses
type OptionResponse struct {
	ID         OptionID   `json:"id"`
	QuestionID QuestionID `json:"question_id"`
	Order      int        `json:"order"`
	Caption    string     `json:"caption"`
	IsCorrect  *bool      `json:"is_correct"`
}

// ToOptionResponse maps an option model to a OptionResponseDTO
func ToOptionResponse(option *Option) *OptionResponse {
	return &OptionResponse{
		ID:         OptionID(option.ID),
		QuestionID: QuestionID(option.QuestionID),
		Order:      option.Order,
		Caption:    option.Caption,
		IsCorrect:  option.IsCorrect,
	}
}

func ToOptionResponses(options *[]Option) *[]OptionResponse {
	optionResponses := make([]OptionResponse, 0)
	for _, option := range *options {
		optionResponses = append(optionResponses, OptionResponse{
			ID:         OptionID(option.ID),
			QuestionID: QuestionID(option.QuestionID),
			Order:      option.Order,
			Caption:    option.Caption,
			IsCorrect:  option.IsCorrect,
		})
	}

	return &optionResponses
}

// ToOptionModel maps a CreateOptionDTO to an Option model
func ToOptionModel(optionDTO *CreateOptionDTO) *Option {
	return &Option{
		QuestionID: optionDTO.QuestionID,
		Order:      optionDTO.Order,
		Caption:    optionDTO.Caption,
		IsCorrect:  optionDTO.IsCorrect,
	}
}

// UpdateOptionModel updates the fields of an Option model from an UpdateOptionDTO
func UpdateOptionModel(option *Option, optionDTO *UpdateOptionDTO) {
	if optionDTO.QuestionID != nil {
		option.QuestionID = *optionDTO.QuestionID
	}
	if optionDTO.Order != nil {
		option.Order = *optionDTO.Order
	}
	if optionDTO.Caption != nil {
		option.Caption = *optionDTO.Caption
	}
	if optionDTO.IsCorrect != nil {
		option.IsCorrect = optionDTO.IsCorrect
	}
}

func (dto *CreateOptionDTO) Validate() error {
	// Check if QuestionID is valid (non-zero)
	if dto.QuestionID == 0 {
		return ErrInvalidQuestionID
	}
	// Check if Order is positive
	if dto.Order <= 0 {
		return ErrInvalidOrder
	}
	// Check if Caption is not empty and within length limits
	if len(strings.TrimSpace(dto.Caption)) == 0 && len(dto.Caption) > 255 {
		return ErrInvalidCaption
	}
	// No validation needed for IsCorrect since it can be nil
	return nil
}

func (dto *UpdateOptionDTO) Validate() error {
	// Ensure at least one field is provided
	if dto.QuestionID == nil && dto.Order == nil && dto.Caption == nil && dto.IsCorrect == nil {
		return ErrAtLeatOneFieldNeededOption
	}
	// Validate each provided field
	if dto.QuestionID != nil && *dto.QuestionID == 0 {
		return ErrInvalidQuestionID
	}
	if dto.Order != nil && *dto.Order <= 0 {
		return ErrInvalidOrder
	}
	if dto.Caption != nil {
		caption := strings.TrimSpace(*dto.Caption)
		if len(caption) == 0 && len(*dto.Caption) > 255 {
			return ErrInvalidCaption
		}
	}
	// No validation needed for IsCorrect since it can be nil
	return nil
}
