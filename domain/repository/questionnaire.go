package domain_repository

import (
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
)

type IQuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *model.Questionnaire) (model.Questionnaire, error)
	GetQuestionnaireByID(questionnaireID model.QuestionnaireID) (model.Questionnaire, error)
	UpdateQuestionaire(questionnaireID model.QuestionnaireID, questionnaire *model.Questionnaire) error
	DeleteQuestionnaire(questionnaireID model.QuestionnaireID) error
	GetUserQuestionnairesCount(userID model.UserID) (int64, error)
	GetQuestionnairesByOwnerID(ownerID model.UserID, page int) ([]Questionnaire, error)
}

type Option struct {
	ID         uint `json:"id"`
	QuestionID uint `json:"question_id"`
	Order      int
	Caption    string
	IsCorrect  *bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Question struct {
	ID                  uint `json:"id"`
	Title               string
	Order               int
	FilePath            *string
	DependsOnQuestionID *uint
	DependsOnOptionID   *uint
	Type                string
	QuestionnaireID     uint `json:"questionnaire_id"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Options             []Option `json:"options" gorm:"foreignKey:QuestionID"`
}

type Questionnaire struct {
	ID                         uint   `json:"id"`
	Title                      string `json:"title"`
	OwnerID                    uint   `json:"owner_id"`
	Status                     string
	CanSubmitFrom              time.Time
	CanSubmitUntil             time.Time
	MaxMinutesToResponse       int
	MaxMinutesToChangeAnswer   int
	MaxMinutesToGivebackAnswer int
	RandomOrSequential         string
	CanBackToPreviousQuestion  bool
	MaxAllowedSubmissionsCount int
	AnswersVisibleFor          string
	CreatedAt                  time.Time
	Questions                  []Question `json:"questions" gorm:"foreignKey:QuestionnaireID"`
}
