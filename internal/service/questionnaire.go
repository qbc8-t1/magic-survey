package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

var (
	ErrQuestionnaireRetrieveFailed = errors.New("failed to retrieve questionnaire")
)

type IQuestionnaireService interface {
	CreateQuestionnaire(questionnaire *model.Questionnaire) (model.Questionnaire, error)
	GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (model.Questionnaire, error)
	UpdateQuestionaire(questionnaireID model.QuestionnaireID, updateData *model.Questionnaire) error
	DeleteQuestionnaire(questionnaireID model.QuestionnaireID) error
	CheckIfUserCanMakeNewQuestionnaire(user model.User) (bool, error)
	// QuestionnairesList() (bool, error)
}

type QuestionnaireService struct {
	questionnaireRepo domain_repository.IQuestionnaireRepository
	userRepo          domain_repository.IUserRepository
}

func (s *QuestionnaireService) GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (model.Questionnaire, error) {
	return s.questionnaireRepo.GetQuestionnaireByID(questionnaireID)
}

func NewQuestionnaireService(questionnaireRepo domain_repository.IQuestionnaireRepository, userRepo domain_repository.IUserRepository) *QuestionnaireService {
	return &QuestionnaireService{
		questionnaireRepo: questionnaireRepo,
		userRepo:          userRepo,
	}
}

func (s *QuestionnaireService) CheckIfUserCanMakeNewQuestionnaire(user model.User) (bool, error) {
	count, err := s.questionnaireRepo.GetUserQuestionnairesCount(model.UserId(user.ID))
	if err != nil {
		return false, errors.Join(errors.New("error in getting questionnaires count"), err)
	}

	if user.MaxQuestionnairesCount == 0 || int(count) < user.MaxQuestionnairesCount {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *QuestionnaireService) CreateQuestionnaire(questionnaire *model.Questionnaire) (model.Questionnaire, error) {
	return s.questionnaireRepo.CreateQuestionnaire(questionnaire)
}

func (s *QuestionnaireService) UpdateQuestionaire(questionnaireID model.QuestionnaireID, updateData *model.Questionnaire) error {
	return s.questionnaireRepo.UpdateQuestionaire(questionnaireID, updateData)
}

func (s *QuestionnaireService) DeleteQuestionnaire(questionnaireID model.QuestionnaireID) error {
	return s.questionnaireRepo.DeleteQuestionnaire(questionnaireID)
}
