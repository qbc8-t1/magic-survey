package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type ICoreService interface {
	Start(questionnaireID model.QuestionnaireID) (*model.QuestionResponse, error)
	Submit(questionID model.QuestionID, answer *model.Answer) error
	Back(model.UserId) (*model.QuestionResponse, error)
	Next(model.UserId) (*model.QuestionResponse, error)
	End() error
}

type CoreService struct {
	submissionRepo    domain_repository.ISubmissionRepository
	questionnaireRepo domain_repository.IQuestionnaireRepository
	questionRepo      domain_repository.IQuestionRepository
	optionRepo        domain_repository.IOptionRepository
}

func NewCoreService(
	submissionRepo domain_repository.ISubmissionRepository,
	questionnaireRepo domain_repository.IQuestionnaireRepository,
	questionRepo domain_repository.IQuestionRepository,
	optionRepo domain_repository.IOptionRepository,
) *CoreService {
	return &CoreService{
		submissionRepo:    submissionRepo,
		questionnaireRepo: questionnaireRepo,
		questionRepo:      questionRepo,
		optionRepo:        optionRepo,
	}
}

func (s *CoreService) Start(questionnaireID model.QuestionnaireID) (*model.QuestionResponse, error) {
	return nil, nil
}

func (s *CoreService) Submit(questionID model.QuestionID, answer *model.Answer) error {
	return nil
}

func (s *CoreService) Back(userID model.UserId) (*model.QuestionResponse, error) {
	// Retrieve the user's active submission
	submission, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active submission found")
		}
		return nil, ErrSubmissionRetrieveFailed
	}

	// Check if LastAnsweredQuestionID is set
	if submission.LastAnsweredQuestionID == nil {
		return nil, errors.New("already at the first question")
	}

	// Get the current question
	currentQuestion, err := s.questionRepo.GetQuestionByID(*submission.LastAnsweredQuestionID)
	if err != nil {
		return nil, ErrQuestionRetrieveFailed
	}

	// Find the previous question
	prevQuestion, err := s.questionnaireRepo.GetPreviousQuestion(submission.QuestionnaireID, currentQuestion.Order)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("already at the first question")
		}
		return nil, ErrQuestionRetrieveFailed
	}

	// Get options if it's a multioption question
	if prevQuestion.Type == model.QuestionsTypeMultioption {
		options, err := s.optionRepo.GetOptionsByQuestionID(prevQuestion.ID)
		if err != nil {
			return nil, ErrOptionRetrieveFailed
		}
		prevQuestion.Options = options
	}

	// Return the previous question
	questionResponse := model.ToQuestionResponse(prevQuestion)
	return questionResponse, nil
}

func (s *CoreService) Next(userID model.UserId) (*model.QuestionResponse, error) {
	// Retrieve the user's active submission
	submission, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active submission found")
		}
		return nil, ErrSubmissionRetrieveFailed
	}

	var nextQuestion *model.Question

	if submission.LastAnsweredQuestionID != nil {
		// Get the current question
		currentQuestion, err := s.questionRepo.GetQuestionByID(*submission.LastAnsweredQuestionID)
		if err != nil {
			return nil, ErrQuestionRetrieveFailed
		}

		// Find the next question
		nextQuestion, err = s.questionnaireRepo.GetNextQuestion(submission.QuestionnaireID, currentQuestion.Order)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("no next question available")
			}
			return nil, ErrQuestionRetrieveFailed
		}
	} else {
		// Get the first question
		nextQuestion, err = s.questionnaireRepo.GetFirstQuestion(submission.QuestionnaireID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("no questions available in this questionnaire")
			}
			return nil, ErrQuestionRetrieveFailed
		}
	}

	// Load options if it's a multioption question
	if nextQuestion.Type == model.QuestionsTypeMultioption {
		options, err := s.optionRepo.GetOptionsByQuestionID(nextQuestion.ID)
		if err != nil {
			return nil, ErrOptionRetrieveFailed
		}
		nextQuestion.Options = options
	}

	// Return the next question
	questionResponse := model.ToQuestionResponse(nextQuestion)
	return questionResponse, nil
}

func (s *CoreService) End() error {
	return nil
}
