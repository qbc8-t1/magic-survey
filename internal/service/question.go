package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

var (
	ErrQuestionIdNotFound = errors.New("question id not found")
)

type IQuestionService interface {
	CreateQuestion(questionDTO *model.CreateQuestionDTO) error
	GetQuestionByID(id uint) (*model.QuestionResponse, error)
	GetQuestionsByID(ids []uint) (*[]model.QuestionResponse, error)
	UpdateQuestion(questionDTO *model.UpdateQuestionDTO) error
	DeleteQuestion(id uint) error
}

type QuestionService struct {
	// dependency injection
	repo domain_repository.IQuestionRepository
}

// NewQuestionService creates a new QuestionService object
func NewQuestionService(repo domain_repository.IQuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) CreateQuestion(questionDTO *model.CreateQuestionDTO) error {
	return nil
}

func (s *QuestionService) GetQuestionByID(id uint) (*model.QuestionResponse, error) {
	question, err := s.repo.GetQuestionByID(id)
	if err != nil {
		return nil, ErrQuestionIdNotFound
	}

	return model.ToQuestionResponse(question), nil
}

func (s *QuestionService) GetQuestionsByID(ids []uint) (*[]model.QuestionResponse, error) {
	return nil, nil
}

func (s *QuestionService) UpdateQuestion(questionDTO *model.UpdateQuestionDTO) error {
	return nil
}

func (s *QuestionService) DeleteQuestion(id uint) error {
	return nil
}
