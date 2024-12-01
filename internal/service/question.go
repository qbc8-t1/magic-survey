package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

var (
	ErrQuestionIdNotFound  = errors.New("question id not found")
	ErrQuestionsIdNotFound = errors.New("questions' ids not found")
	ErrQuestionOnCreate    = errors.New("cannot create the question")
	ErrQuestionNotFound    = errors.New("question not found")
	ErrQuestionOnUpdate    = errors.New("cannot update the question")
	ErrQuestionOnDelete    = errors.New("cannot delete the question")
)

type IQuestionService interface {
	CreateQuestion(questionDTO *model.CreateQuestionDTO) error
	GetQuestionByID(id uint) (*model.QuestionResponse, error)
	GetAllQuestions() (*[]model.QuestionResponse, error)
	UpdateQuestion(id uint, questionDTO *model.UpdateQuestionDTO) error
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
	question := model.ToQuestionModel(questionDTO)
	err := s.repo.CreateQuestion(question)
	if err != nil {
		return ErrQuestionOnCreate
	}

	return nil
}

func (s *QuestionService) GetQuestionByID(id uint) (*model.QuestionResponse, error) {
	question, err := s.repo.GetQuestionByID(id)
	if err != nil {
		return nil, ErrQuestionIdNotFound
	}

	return model.ToQuestionResponse(question), nil
}

func (s *QuestionService) GetAllQuestions() (*[]model.QuestionResponse, error) {
	questions, err := s.repo.GetAllQuestions()
	if err != nil {
		return nil, ErrQuestionsIdNotFound
	}

	return model.ToQuestionResponses(questions), nil
}

func (s *QuestionService) UpdateQuestion(id uint, questionDTO *model.UpdateQuestionDTO) error {
	// Check if the question exists
	existingQuestion, err := s.repo.GetQuestionByID(id)
	if err != nil {
		return ErrQuestionNotFound
	}

	// Map changes to the existing question
	model.UpdateQuestionModel(existingQuestion, questionDTO)

	// Persist the changes
	err = s.repo.UpdateQuestion(existingQuestion)
	if err != nil {
		return ErrQuestionOnUpdate
	}

	return nil
}

func (s *QuestionService) DeleteQuestion(id uint) error {
	// Check if the question exists
	_, err := s.repo.GetQuestionByID(id)
	if err != nil {
		return ErrQuestionNotFound
	}

	// Delete the question
	err = s.repo.DeleteQuestion(id)
	if err != nil {
		return ErrQuestionOnDelete
	}

	return nil
}
