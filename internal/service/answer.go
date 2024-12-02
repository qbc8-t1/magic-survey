package service

import (
	"errors"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

// Custom errors
var (
	// General errors
	ErrAnswerNotFound      = errors.New("answer not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrSubmissionNotFound  = errors.New("submission not found")
	ErrOptionNotFound      = errors.New("option not found")
	ErrInvalidAnswerID     = errors.New("invalid answer ID")
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrInvalidSubmissionID = errors.New("invalid submission ID")
	ErrInvalidOptionID     = errors.New("invalid option ID")

	// Creation errors
	ErrAnswerCreateFailed = errors.New("failed to create answer")

	// Update errors
	ErrAnswerUpdateFailed = errors.New("failed to update answer")

	// Delete errors
	ErrAnswerDeleteFailed = errors.New("failed to delete answer")

	// Retrieval errors
	ErrAnswerRetrieveFailed     = errors.New("failed to retrieve answer")
	ErrUserRetrieveFailed       = errors.New("failed to retrieve user")
	ErrSubmissionRetrieveFailed = errors.New("failed to retrieve submission")
	ErrOptionRetrieveFailed     = errors.New("failed to retrieve option")
)

type IAnswerService interface {
	CreateAnswer(AnswerDTO *model.CreateAnswerDTO) error
	GetAnswerByID(id model.AnswerID) (*model.AnswerResponse, error)
	UpdateAnswer(id model.AnswerID, AnswerDTO *model.UpdateAnswerDTO) error
	DeleteAnswer(id model.AnswerID) error
}

type AnswerService struct {
	answerRepo     domain_repository.IAnswerRepository
	userRepo       domain_repository.IUserRepository
	submissionRepo domain_repository.ISubmissionRepository
	questionRepo   domain_repository.IQuestionRepository
	optionRepo     domain_repository.IOptionRepository
}

func NewAnswerService(
	answerRepo domain_repository.IAnswerRepository,
	userRepo domain_repository.IUserRepository,
	submissionRepo domain_repository.ISubmissionRepository,
	questionRepo domain_repository.IQuestionRepository,
	optionRepo domain_repository.IOptionRepository,
) *AnswerService {
	return &AnswerService{
		answerRepo:     answerRepo,
		userRepo:       userRepo,
		submissionRepo: submissionRepo,
		questionRepo:   questionRepo,
		optionRepo:     optionRepo,
	}
}

func (s *AnswerService) CreateAnswer(answerDTO *model.CreateAnswerDTO) error {
	// Validate UserID
	if answerDTO.UserID == 0 {
		return ErrInvalidUserID
	}
	// Check if User exists
	_, err := s.userRepo.GetUserByID(model.UserId(answerDTO.UserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return ErrUserRetrieveFailed
	}

	// Validate SubmissionID
	if answerDTO.SubmissionID == 0 {
		return ErrInvalidSubmissionID
	}
	// Check if Submission exists
	_, err = s.submissionRepo.GetSubmissionByID(answerDTO.SubmissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubmissionNotFound
		}
		return ErrSubmissionRetrieveFailed
	}

	// Validate QuestionID
	if answerDTO.QuestionID == 0 {
		return ErrInvalidQuestionID
	}
	// Check if Question exists
	_, err = s.questionRepo.GetQuestionByID(model.QuestionID(answerDTO.QuestionID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return ErrQuestionRetrieveFailed
	}

	// Validate OptionID if provided
	if answerDTO.OptionID != nil {
		if *answerDTO.OptionID == 0 {
			return ErrInvalidOptionID
		}
		// Check if Option exists
		_, err = s.optionRepo.GetOptionByID(*answerDTO.OptionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOptionNotFound
			}
			return ErrOptionRetrieveFailed
		}
	}

	// Convert DTO to model
	answer := model.ToAnswerModel(answerDTO)
	answer.CreatedAt = time.Now()

	// Create the answer
	err = s.answerRepo.CreateAnswer(answer)
	if err != nil {
		return ErrAnswerCreateFailed
	}

	return nil
}

func (s *AnswerService) GetAnswerByID(id model.AnswerID) (*model.AnswerResponse, error) {
	Answer, err := s.answerRepo.GetAnswerByID(id)
	if err != nil {
		return nil, ErrAnswerNotFound
	}

	return model.ToAnswerResponse(Answer), nil
}

func (s *AnswerService) UpdateAnswer(id model.AnswerID, answerDTO *model.UpdateAnswerDTO) error {
	// Validate Answer ID
	if id == 0 {
		return ErrInvalidAnswerID
	}

	// Check if the Answer exists
	existingAnswer, err := s.answerRepo.GetAnswerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAnswerNotFound
		}
		return ErrAnswerRetrieveFailed
	}

	// If OptionID is being updated
	if answerDTO.OptionID != nil {
		if *answerDTO.OptionID == 0 {
			return ErrInvalidOptionID
		}
		// Check if Option exists
		_, err = s.optionRepo.GetOptionByID(*answerDTO.OptionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrOptionNotFound
			}
			return ErrOptionRetrieveFailed
		}
	}

	existingAnswer.UpdatedAt = time.Now()
	model.UpdateAnswerModel(existingAnswer, answerDTO)

	// Persist the changes
	err = s.answerRepo.UpdateAnswer(existingAnswer)
	if err != nil {
		return ErrAnswerUpdateFailed
	}

	return nil
}

func (s *AnswerService) DeleteAnswer(id model.AnswerID) error {
	// Check if the Answer exists
	_, err := s.answerRepo.GetAnswerByID(id)
	if err != nil {
		return ErrAnswerNotFound
	}

	// Delete the Answer
	err = s.answerRepo.DeleteAnswer(id)
	if err != nil {
		return ErrAnswerDeleteFailed
	}

	return nil
}
