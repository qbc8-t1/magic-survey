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
	ErrQuestionNotFound      = errors.New("question not found")
	ErrQuestionnaireNotFound = errors.New("questionnaire not found")
	ErrInvalidQuestionID     = errors.New("invalid question ID")

	// Creation errors
	ErrQuestionCreateFailed = errors.New("failed to create question")

	// Update errors
	ErrQuestionUpdateFailed = errors.New("failed to update question")

	// Delete errors
	ErrQuestionDeleteFailed = errors.New("failed to delete question")

	// Retrieval errors
	ErrQuestionRetrieveFailed = errors.New("failed to retrieve question")
)

type IQuestionService interface {
	CreateQuestion(questionDTO *model.CreateQuestionDTO) error
	GetQuestionByID(id model.QuestionID) (*model.QuestionResponse, error)
	GetQuestionsByQuestionnaireID(quesionnaireID model.QuestionnaireID) (*[]model.QuestionResponse, error)
	UpdateQuestion(id model.QuestionID, questionDTO *model.UpdateQuestionDTO) error
	DeleteQuestion(id model.QuestionID) error
}

type QuestionService struct {
	// dependency injection
	questionRepo      domain_repository.IQuestionRepository
	questionnaireRepo domain_repository.IQuestionnaireRepo
}

// NewQuestionService creates a new QuestionService object
func NewQuestionService(questionRepo domain_repository.IQuestionRepository,
	questionnaireRepo domain_repository.IQuestionnaireRepo) *QuestionService {
	return &QuestionService{
		questionRepo:      questionRepo,
		questionnaireRepo: questionnaireRepo,
	}
}

func (s *QuestionService) CreateQuestion(questionDTO *model.CreateQuestionDTO) error {
	// Check if the Questionnaire exists
	questionnaireID := questionDTO.QuestionnaireID
	_, err := s.questionnaireRepo.GetQuestionnaireByID(questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionnaireNotFound
		}
		return ErrQuestionnaireRetrieveFailed
	}

	// Convert DTO to model
	question := model.ToQuestionModel(questionDTO)
	question.CreatedAt = time.Now()

	// Create the question
	err = s.questionRepo.CreateQuestion(question)
	if err != nil {
		return ErrQuestionCreateFailed
	}

	return nil
}

func (s *QuestionService) UpdateQuestion(id model.QuestionID, questionDTO *model.UpdateQuestionDTO) error {
	// Check if the question exists
	existingQuestion, err := s.questionRepo.GetQuestionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return ErrQuestionRetrieveFailed
	}

	// If QuestionnaireID is being updated, check if the new Questionnaire exists
	if questionDTO.QuestionnaireID != nil {
		_, err := s.questionnaireRepo.GetQuestionnaireByID(*questionDTO.QuestionnaireID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrQuestionnaireNotFound
			}
			return ErrQuestionnaireRetrieveFailed
		}
	}

	existingQuestion.UpdatedAt = time.Now()
	model.UpdateQuestionModel(existingQuestion, questionDTO)

	// Persist the changes
	err = s.questionRepo.UpdateQuestion(existingQuestion)
	if err != nil {
		return ErrQuestionUpdateFailed
	}

	return nil
}

func (s *QuestionService) GetQuestionByID(id model.QuestionID) (*model.QuestionResponse, error) {
	question, err := s.questionRepo.GetQuestionByID(id)
	if err != nil {
		return nil, ErrQuestionNotFound
	}

	return model.ToQuestionResponse(question), nil
}

func (s *QuestionService) GetQuestionsByQuestionnaireID(questionnaireID model.QuestionnaireID) (*[]model.QuestionResponse, error) {
	questions, err := s.questionRepo.GetQuestionsByQuestionnaireID(questionnaireID)
	if err != nil {
		return nil, ErrQuestionNotFound
	}

	return model.ToQuestionResponses(questions), nil
}

func (s *QuestionService) DeleteQuestion(id model.QuestionID) error {
	// Check if the question exists
	_, err := s.questionRepo.GetQuestionByID(id)
	if err != nil {
		return ErrQuestionNotFound
	}

	// Delete the question
	err = s.questionRepo.DeleteQuestion(id)
	if err != nil {
		return ErrQuestionDeleteFailed
	}

	return nil
}
