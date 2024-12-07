package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

// CRUD

// createQuestionnaire()

// GetQuestionnaire()

// DeleteQuestionnaire()

// UpdateQuestionaire()

type IQuestionnaireService interface {
	CreateQuestionnaire(questionnaire *model.Questionnaire) (model.Questionnaire, error)
	GetQuestionnaireByID(id uint) (*model.Questionnaire, error)
	UpdateQuestionaire(questionnaire *model.Questionnaire) error
	DeleteQuestionnaire(id uint) error
	CheckIfUserCanMakeNewQuestionnaire(user model.User) (bool, error)
}

type QuestionnaireService struct {
	questionnaireRepo domain_repository.IQuestionnaireRepository
	userRepo          domain_repository.IUserRepository
}

func NewQuestionnaireService(questionnaireRepo domain_repository.IQuestionnaireRepository, userRepo domain_repository.IUserRepository) IQuestionnaireService {
	return &QuestionnaireService{
		questionnaireRepo: questionnaireRepo,
		userRepo:          userRepo,
	}
}

func (s *QuestionnaireService) CheckIfUserCanMakeNewQuestionnaire(user model.User) (bool, error) {
	count, err := s.questionnaireRepo.GetUserQuestionnairesCount(user.ID)
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

func (s *QuestionnaireService) GetQuestionnaireByID(id uint) (*model.Questionnaire, error) {
	return nil, nil
}

func (s *QuestionnaireService) UpdateQuestionaire(questionnaire *model.Questionnaire) error {
	return s.questionnaireRepo.UpdateQuestionaire(questionnaire)
}

func (s *QuestionnaireService) DeleteQuestionnaire(id uint) error {
	return nil
}
