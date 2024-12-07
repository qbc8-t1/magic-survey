package service

import (
	"errors"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

// TODO:
// func GetNextQuestionWithDependency(nextQuestion *model.Question, s *CoreService, submission *model.Submission) (*model.Question, error) {
// 	if nextQuestion.DependsOnQuestionID != nil {
// 		answer, err := s.AnswerRepo.GetAnswerBySubmissionIDAndQuestionID(submission.ID, model.QuestionID(*nextQuestion.DependsOnQuestionID))
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return nil, ErrQuestionNotFound
// 			}

// 			return nil, ErrAnswerRetrieveFailed
// 		}
// 		if nextQuestion.Type == model.QuestionsTypeDescriptive && answer.AnswerText == nil {
// 			return s.QuestionRepo.GetQuestionByID(model.QuestionID(*nextQuestion.DependsOnQuestionID))
// 		} else if nextQuestion.Type == model.QuestionsTypeDescriptive && answer.AnswerText != nil {
// 			return nextQuestion, nil
// 		} else if nextQuestion.Type == model.QuestionsTypeMultioption {
// 			dq, err := s.QuestionRepo.GetQuestionByID(model.QuestionID(*nextQuestion.DependsOnQuestionID))
// 			// error handling
// 			answer_dq, err := s.AnswerRepo.GetAnswerBySubmissionIDAndQuestionID(submission.ID, dq.ID)
// 			// TODO: error handling
// 			if answer_dq.OptionID == *model.OptionID(nextQuestion.DependsOnOptionID) {
// 				return nextQuestion, nil
// 			} else {
// 				submission.LastAnsweredQuestionID = &nextQuestion.ID
// 				err = s.SubmissionRepo.UpdateSubmission(submission)
// 				// error handling
// 				return GetNextQuestionWithDependency(dq, s, submission)
// 			}
// 		}
// 	}
// 	return nextQuestion, nil
// }

type ICoreService interface {
	Start(questionnaireID model.QuestionnaireID, userID model.UserId) (*model.QuestionResponse, error)
	Submit(questionID model.QuestionID, answer *model.Answer, userID model.UserId) error
	Back(userID model.UserId) (*model.QuestionResponse, error)
	Next(userID model.UserId) (*model.QuestionResponse, error)
	End(userID model.UserId) error
}

type CoreService struct {
	submissionRepo    domain_repository.ISubmissionRepository
	questionnaireRepo domain_repository.IQuestionnaireRepository
	questionRepo      domain_repository.IQuestionRepository
	optionRepo        domain_repository.IOptionRepository
	answerRepo        domain_repository.IAnswerRepository
}

func NewCoreService(
	submissionRepo domain_repository.ISubmissionRepository,
	questionnaireRepo domain_repository.IQuestionnaireRepository,
	questionRepo domain_repository.IQuestionRepository,
	optionRepo domain_repository.IOptionRepository,
	answerRepo domain_repository.IAnswerRepository) *CoreService {
	return &CoreService{
		submissionRepo:    submissionRepo,
		questionnaireRepo: questionnaireRepo,
		questionRepo:      questionRepo,
		optionRepo:        optionRepo,
		answerRepo:        answerRepo,
	}
}

func (s *CoreService) Start(questionnaireID model.QuestionnaireID, userID model.UserId) (*model.QuestionResponse, error) {
	// Check if questionnaire exists
	_, err := s.questionnaireRepo.GetQuestionnaireByID(questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrorQuestionnaireNotFound

		}
		return nil, ErrQuestionnaireRetrieveFailed
	}

	// Check if user already has an active submission for this questionnaire
	activeSubmission, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSubmissionRetrieveFailed
	}

	// If there's an active submission for a different questionnaire, we might need to handle that case
	// For now, let's assume the user can only have one active submission at a time.
	if activeSubmission != nil && activeSubmission.QuestionnaireID == questionnaireID {
		// If the submission is for the same questionnaire, resume from where they left off
		// Get the next question based on LastAnsweredQuestionID
		var nextQuestion *model.Question
		if activeSubmission.LastAnsweredQuestionID != nil {
			currentQ, err := s.questionRepo.GetQuestionByID(*activeSubmission.LastAnsweredQuestionID)
			if err != nil {
				return nil, ErrQuestionRetrieveFailed
			}
			nextQuestion, err = s.questionnaireRepo.GetNextQuestion(questionnaireID, currentQ.Order)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrNoNextQuestionAvailable
				}
				return nil, ErrQuestionRetrieveFailed
			}
		} else {
			// If LastAnsweredQuestionID is nil, this means the user has not answered anything yet
			nextQuestion, err = s.questionnaireRepo.GetFirstQuestion(questionnaireID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrNoQuestionsInQuestionnaire
				}
				return nil, ErrQuestionRetrieveFailed
			}
		}

		if nextQuestion.Type == model.QuestionsTypeMultioption {
			options, err := s.optionRepo.GetOptionsByQuestionID(nextQuestion.ID)
			if err != nil {
				return nil, ErrOptionRetrieveFailed
			}
			nextQuestion.Options = options
		}

		return model.ToQuestionResponse(nextQuestion), nil
	}

	// Otherwise, create a new submission
	newSubmission := &model.Submission{
		QuestionnaireID: questionnaireID,
		UserID:          userID,
		Status:          model.SubmissionsStatusAnswering,
		// LastAnsweredQuestionID remains nil at start
	}

	err = s.submissionRepo.CreateSubmission(newSubmission)
	if err != nil {
		return nil, ErrSubmissionCreateFailed
	}

	// Get the first question of the questionnaire
	firstQuestion, err := s.questionnaireRepo.GetFirstQuestion(questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoQuestionsInQuestionnaire
		}
		return nil, ErrQuestionRetrieveFailed
	}

	// Load options if multioption
	if firstQuestion.Type == model.QuestionsTypeMultioption {
		options, err := s.optionRepo.GetOptionsByQuestionID(firstQuestion.ID)
		if err != nil {
			return nil, ErrOptionRetrieveFailed
		}
		firstQuestion.Options = options
	}

	return model.ToQuestionResponse(firstQuestion), nil
}

func (s *CoreService) Submit(questionID model.QuestionID, answer *model.Answer, userID model.UserId) error {
	// Retrieve the user's active submission
	submission, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoActiveSubmissionFound
		}
		return ErrSubmissionRetrieveFailed
	}

	// Ensure that the question exists and belongs to the questionnaire
	question, err := s.questionRepo.GetQuestionByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrQuestionNotFound
		}
		return ErrQuestionRetrieveFailed
	}

	// Optional: Validate that the question belongs to the same questionnaire
	if question.QuestionnaireID != submission.QuestionnaireID {
		return ErrQuestionDoesNotBelongToQuestionnaire
	}

	// Prepare the answer for submission
	answer.SubmissionID = submission.ID
	answer.QuestionID = questionID

	// Check if an answer already exists for this question
	existingAnswer, err := s.answerRepo.GetAnswerBySubmissionIDAndQuestionID(submission.ID, questionID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAnswerRetrieveFailed
		}
		// If not found, we create a new answer
		err = s.answerRepo.CreateAnswer(answer)
		if err != nil {
			return ErrAnswerCreateFailed
		}
	} else {
		// If found, we update the existing answer
		existingAnswer.AnswerText = answer.AnswerText
		existingAnswer.OptionID = answer.OptionID
		err = s.answerRepo.UpdateAnswer(existingAnswer)
		if err != nil {
			return ErrAnswerUpdateFailed
		}
	}

	// Update the LastAnsweredQuestionID in the submission
	submission.LastAnsweredQuestionID = &questionID
	err = s.submissionRepo.UpdateSubmission(submission)
	if err != nil {
		return ErrSubmissionUpdateFailed
	}

	return nil
}

// TODO:
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

// TODO:
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

	// if nextQuestion.DependsOnQuestionID != nil {
	// 	answer, err := s.answerRepo.GetAnswerBySubmissionIDAndQuestionID(submission.ID, model.QuestionID(*nextQuestion.DependsOnQuestionID))
	// 	if err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return nil, errors.New("next question not found")
	// 		}

	// 		return nil, ErrAnswerRetrieveFailed
	// 	}
	// 	if nextQuestion.Type == model.QuestionsTypeDescriptive && answer.AnswerText == nil {
	// 		nextQuestion, err = s.questionRepo.GetQuestionByID(model.QuestionID(*nextQuestion.DependsOnQuestionID))
	// 	} else if nextQuestion.Type == model.QuestionsTypeDescriptive && answer.AnswerText == nil {

	// 	} else if nextQuestion.Type == model.QuestionsTypeMultioption {
	// 		answer_dq, err := s.answerRepo.GetAnswerBySubmissionIDAndQuestionID(submission.ID, model.QuestionID(*nextQuestion.DependsOnQuestionID))
	// 		// TODO: error handling
	// 		if answer_dq.OptionID != (*model.OptionID)(nextQuestion.DependsOnOptionID) {
	// 			nextQuestion = s.
	// 		}
	// 	}
	// }

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

func (s *CoreService) End(userID model.UserId) error {
	// Retrieve the user's active submission
	submission, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoActiveSubmissionFound
		}
		return ErrSubmissionRetrieveFailed
	}

	// Mark the submission as submitted
	now := time.Now()
	submission.Status = model.SubmissionsStatusSubmitted
	submission.SubmittedAt = &now

	err = s.submissionRepo.UpdateSubmission(submission)
	if err != nil {
		return ErrSubmissionUpdateFailed
	}

	return nil
}
