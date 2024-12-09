package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

var (
	ErrNoAttemptsLeft          = errors.New("no attempts left for this questionnaire")
	ErrNoActiveSubmission      = errors.New("no active submission found for user")
	ErrWrongQuestionOrder      = errors.New("attempt to submit a question not matching current question")
	ErrTimeExpired             = errors.New("time to answer expired")
	ErrCannotGoBack            = errors.New("cannot go back to previous question")
	ErrNoNextQuestion          = errors.New("no next question available")
	ErrSubmissionNotAnswering  = errors.New("submission not in answering status")
	ErrAlreadySubmitted        = errors.New("submission is already submitted")
	ErrDependencyNotSatisfied  = errors.New("cannot answer this question due to unsatisfied dependencies")
	ErrIncompleteQuestionnaire = errors.New("cannot end submission, not all questions answered")
	ErrInvalidAnswer           = errors.New("invalid answer format for question type")
)

// ICoreService defines the main operations of the questionnaire
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

// Start begins a new submission if allowed
func (s *CoreService) Start(questionnaireID model.QuestionnaireID, userID model.UserId) (*model.QuestionResponse, error) {
	questionnaire, err := s.questionnaireRepo.GetQuestionnaireByID(questionnaireID)
	if err != nil || questionnaire == nil || questionnaire.Status != model.QuestionnaireStatusOpen {
		return nil, ErrQuestionnaireNotFound
	}

	// Check attempts left
	subCount := 0
	for _, sub := range questionnaire.Submissions {
		if sub.UserID == userID && (sub.Status == model.SubmissionsStatusSubmitted || sub.Status == model.SubmissionsStatusAnswering) {
			subCount++
		}
	}
	if subCount >= questionnaire.MaxAllowedSubmissionsCount {
		return nil, ErrNoAttemptsLeft
	}

	// If user currently has an active submission for this questionnaire, either continue that or fail
	activeSub, _ := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if activeSub != nil && activeSub.QuestionnaireID == questionnaireID {
		// They already started this questionnaire, return current question if still answering
		if activeSub.Status == model.SubmissionsStatusAnswering {
			// fetch current question
			q, err := s.questionRepo.GetQuestionByID(*activeSub.CurrentQuestionID)
			if err != nil {
				return nil, err
			}
			return model.ToQuestionResponse(q), nil
		}
		// If status is not answering, they need to start a new submission if possible
	}

	// Create a new submission
	questions, err := s.questionRepo.GetQuestionsByQuestionnaireID(questionnaireID)
	if err != nil || len(*questions) == 0 {
		return nil, errors.New("no questions found for questionnaire")
	}

	var questionOrder []model.QuestionID

	if questionnaire.RandomOrSequential == model.QuestionnaireTypeRandom {
		questions = RandomizeQuestions(*questions)
	}
	for _, q := range *questions {
		questionOrder = append(questionOrder, q.ID)
	}

	// Create new submission
	firstQuestionID := questionOrder[0]

	newSub := &model.Submission{
		QuestionnaireID:   questionnaireID,
		UserID:            userID,
		Status:            model.SubmissionsStatusAnswering,
		CurrentQuestionID: &firstQuestionID,
		QuestionOrder:     questionOrder,
	}

	err = s.submissionRepo.CreateSubmission(newSub)
	if err != nil {
		return nil, err
	}

	// Return first question
	firstQuestion := (*questions)[0]
	return model.ToQuestionResponse(&firstQuestion), nil
}

func (s *CoreService) Submit(questionID model.QuestionID, answer *model.Answer, userID model.UserId) error {
	sub, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil || sub == nil {
		return ErrNoActiveSubmission
	}
	if sub.Status != model.SubmissionsStatusAnswering {
		return ErrSubmissionNotAnswering
	}

	// Check time limit
	if isTimeExpired(sub) {
		finalizeSubmission(sub)
		return s.submissionRepo.UpdateSubmission(sub)
	}

	// Ensure current question matches questionID
	if sub.CurrentQuestionID == nil || *sub.CurrentQuestionID != questionID {
		return ErrWrongQuestionOrder
	}

	// Check question type and validate answer format
	q, err := s.questionRepo.GetQuestionByID(questionID)
	if err != nil {
		return err
	}

	// Check dependencies
	err = s.checkDependencies(*sub, questionID)
	if err != nil {
		return err
	}

	switch q.Type {
	case model.QuestionsTypeMultioption:
		// Multioption questions must have OptionID set and AnswerText not set
		if answer.OptionID == nil || (answer.AnswerText != nil && *answer.AnswerText != "") {
			return ErrInvalidAnswer
		}
	case model.QuestionsTypeDescriptive:
		// Descriptive questions must have AnswerText set and OptionID not set
		if answer.AnswerText == nil || (answer.OptionID != nil) {
			return ErrInvalidAnswer
		}
	default:
		return ErrInvalidAnswer
	}

	// Save answer
	answer.UserID = userID
	answer.SubmissionID = sub.ID
	answer.QuestionID = questionID

	err = s.answerRepo.CreateAnswer(answer)
	if err != nil {
		return err
	}

	// Update submission last answered question
	sub.LastAnsweredQuestionID = &questionID

	return s.submissionRepo.UpdateSubmission(sub)
}

func (s *CoreService) Back(userID model.UserId) (*model.QuestionResponse, error) {
	sub, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil || sub == nil {
		return nil, ErrNoActiveSubmission
	}
	if sub.Status != model.SubmissionsStatusAnswering {
		return nil, ErrSubmissionNotAnswering
	}

	// Check time limit
	if isTimeExpired(sub) {
		finalizeSubmission(sub)
		s.submissionRepo.UpdateSubmission(sub)
		return nil, ErrTimeExpired
	}

	// Check if questionnaire allows going back
	qnr, err := s.questionnaireRepo.GetQuestionnaireByID(sub.QuestionnaireID)
	if err != nil || qnr == nil {
		return nil, ErrQuestionnaireNotFound
	}
	if !qnr.CanBackToPreviousQuestion {
		return nil, ErrCannotGoBack
	}

	// Find previous question by order
	currentQ, err := s.questionRepo.GetQuestionByID(*sub.CurrentQuestionID)
	if err != nil {
		return nil, ErrQuestionNotFound
	}

	prevQ, err := s.questionnaireRepo.GetPreviousQuestion(sub.QuestionnaireID, currentQ.Order)
	if err != nil || prevQ == nil {
		return nil, ErrCannotGoBack
	}

	sub.CurrentQuestionID = &prevQ.ID
	err = s.submissionRepo.UpdateSubmission(sub)
	if err != nil {
		return nil, err
	}
	return model.ToQuestionResponse(prevQ), nil
}

func (s *CoreService) Next(userID model.UserId) (*model.QuestionResponse, error) {
	sub, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil || sub == nil {
		return nil, ErrNoActiveSubmission
	}
	if sub.Status != model.SubmissionsStatusAnswering {
		return nil, ErrSubmissionNotAnswering
	}

	// Check time limit
	if isTimeExpired(sub) {
		finalizeSubmission(sub)
		s.submissionRepo.UpdateSubmission(sub)
		return nil, ErrTimeExpired
	}

	// Ensure current question is answered
	if sub.CurrentQuestionID == nil {
		return nil, ErrNoNextQuestion // No current question means we are done
	}

	_, err = s.answerRepo.GetAnswerBySubmissionIDAndQuestionID(sub.ID, *sub.CurrentQuestionID)
	if err != nil {
		// Current question not answered yet
		return nil, errors.New("please submit the current question before moving to the next")
	}

	// Find the next question based on QuestionOrder
	var nextQuestionID *model.QuestionID
	for i, qID := range sub.QuestionOrder {
		if qID == *sub.CurrentQuestionID && i < len(sub.QuestionOrder)-1 {
			nextQuestionID = &sub.QuestionOrder[i+1]
			break
		}
	}

	// Handle no next question
	if nextQuestionID == nil {
		sub.CurrentQuestionID = nil
		err = s.submissionRepo.UpdateSubmission(sub)
		if err != nil {
			return nil, err
		}
		return nil, ErrNoNextQuestion
	}

	// Fetch the next question
	nextQ, err := s.questionRepo.GetQuestionByID(*nextQuestionID)
	if err != nil {
		return nil, err
	}

	// Update the submission with the new current question
	sub.CurrentQuestionID = nextQuestionID
	err = s.submissionRepo.UpdateSubmission(sub)
	if err != nil {
		return nil, err
	}

	return model.ToQuestionResponse(nextQ), nil
}

func (s *CoreService) End(userID model.UserId) error {
	sub, err := s.submissionRepo.GetActiveSubmissionByUserID(userID)
	if err != nil || sub == nil {
		return ErrNoActiveSubmission
	}

	if sub.Status != model.SubmissionsStatusAnswering {
		return ErrSubmissionNotAnswering
	}

	// Check time
	if isTimeExpired(sub) {
		// Notify user that time is over
		err = s.safelyFinalizeSubmission(sub)
		if err != nil {
			return err
		}
		return errors.New("time is over, submission has been finalized")
	}

	// Ensure all questions have been answered before ending
	unanswered, err := s.questionRepo.GetUnansweredQuestions(sub.QuestionnaireID, sub.ID)
	if err != nil {
		return err
	}
	if unanswered != nil && len(*unanswered) > 0 {
		return ErrIncompleteQuestionnaire
	}

	// All answered, finalize
	err = s.safelyFinalizeSubmission(sub)
	if err != nil {
		return err
	}
	return nil
}

// Helper Methods

func isTimeExpired(sub *model.Submission) bool {
	// Retrieve questionnaire to check max time
	qnr := sub.Questionnaire
	if qnr.MaxMinutesToResponse <= 0 {
		// no time limit
		return false
	}
	if sub.CreatedAt.IsZero() {
		return false
	}
	elapsed := time.Since(sub.CreatedAt)
	maxDuration := time.Duration(qnr.MaxMinutesToResponse) * time.Minute
	return elapsed > maxDuration
}

func finalizeSubmission(sub *model.Submission) {
	sub.Status = model.SubmissionsStatusSubmitted
	now := time.Now()
	sub.SubmittedAt = &now
	elapsed := time.Since(sub.CreatedAt)
	subMin := int(elapsed.Minutes())
	sub.SpentMinutes = &subMin
	sub.CurrentQuestionID = nil
}

// Helper Method: Finalize Submission Safely
func (s *CoreService) safelyFinalizeSubmission(sub *model.Submission) error {
	finalizeSubmission(sub)
	return s.submissionRepo.UpdateSubmission(sub)
}

func (s *CoreService) getNextValidQuestion(sub model.Submission) (*model.Question, error) {
	if sub.CurrentQuestionID == nil {
		return nil, ErrNoNextQuestion
	}

	currentQ, err := s.questionRepo.GetQuestionByID(*sub.CurrentQuestionID)
	if err != nil {
		return nil, err
	}

	// Get next question by order
	nextQ, err := s.questionnaireRepo.GetNextQuestion(sub.QuestionnaireID, currentQ.Order)
	if err != nil || nextQ == nil {
		return nil, ErrNoNextQuestion
	}

	// Check dependencies for next question and skip if not satisfied
	for nextQ != nil {
		if s.isDependencySatisfied(sub, nextQ) {
			return nextQ, nil
		}
		// Get next one
		nextQ, err = s.questionnaireRepo.GetNextQuestion(sub.QuestionnaireID, nextQ.Order)
		if err != nil || nextQ == nil {
			return nil, ErrNoNextQuestion
		}
	}

	return nil, ErrNoNextQuestion
}

func (s *CoreService) isDependencySatisfied(sub model.Submission, q *model.Question) bool {
	// If no dependency, it's satisfied
	if q.DependsOnQuestionID == nil {
		return true
	}

	// If depends on a previous question
	ans, err := s.answerRepo.GetAnswerBySubmissionIDAndQuestionID(sub.ID, *q.DependsOnQuestionID)
	if err != nil || ans == nil {
		// no answer for dependency question means not satisfied
		return false
	}

	// If also depends on an option, check if answered option matches
	if q.DependsOnOptionID != nil {
		// If the answered option is not the one required, dependency fails
		if ans.OptionID == nil || *ans.OptionID != *q.DependsOnOptionID {
			return false
		}
	}

	return true
}

func (s *CoreService) checkDependencies(sub model.Submission, questionID model.QuestionID) error {
	q, err := s.questionRepo.GetQuestionByID(questionID)
	if err != nil {
		return err
	}
	if q.DependsOnQuestionID == nil {
		return nil
	}
	// depends on a previous question
	ans, err := s.answerRepo.GetAnswerBySubmissionIDAndQuestionID(sub.ID, *q.DependsOnQuestionID)
	if err != nil {
		return ErrDependencyNotSatisfied
	}

	if q.DependsOnOptionID != nil {
		if ans.OptionID == nil || *ans.OptionID != *q.DependsOnOptionID {
			return ErrDependencyNotSatisfied
		}
	}
	return nil
}

// RandomizeQuestions generates a randomized order of questions for the questionnaire
func RandomizeQuestions(questions []model.Question) *[]model.Question {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random generator
	r.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
	return &questions
}
