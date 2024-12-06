package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

type IOptionRepository interface {
	CreateOption(option *model.Option) error
	GetOptionByID(id model.OptionID) (*model.Option, error)
	GetOptionsByQuestionID(questionID model.QuestionID) (*[]model.Option, error)
	UpdateOption(option *model.Option) error
	DeleteOption(id model.OptionID) error
}
