package service

import (
	"errors"

	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type IQuestionService interface {
}

type QuestionService struct {
	// dependency injection
	repo domain_repository.IQuestionRepository
}

// NewQuestionService creates a new QuestionService object
func NewQuestionService(repo domain_repository.IQuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) IsQuestionForQuestionnaire(questionID uint, questionnaireID uint) (bool, error) {
	question, err := s.repo.FindQuestionByQuestionIDAndQuestionnaireID(questionID, questionnaireID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if question.ID == 0 {
		return false, nil
	}
	return true, nil
}
