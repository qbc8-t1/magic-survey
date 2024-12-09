package repository

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) domain_repository.IAnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) CreateAnswer(answer *model.Answer) error {
	return r.db.Create(&answer).Error
}

func (r *AnswerRepository) GetAnswerByID(id model.AnswerID) (*model.Answer, error) {
	var answer model.Answer
	result := r.db.Preload("Option").Preload("Question.Options").Preload("Submission.Questionnaire").First(&answer, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &answer, nil
}

func (r *AnswerRepository) UpdateAnswer(answer *model.Answer) error {
	return r.db.Save(&answer).Error
}

func (r *AnswerRepository) DeleteAnswer(id model.AnswerID) error {
	return r.db.Delete(&model.Answer{}, id).Error
}

func (r *AnswerRepository) GetAnswerBySubmissionIDAndQuestionID(submissionID model.SubmissionID, questionID model.QuestionID) (*model.Answer, error) {
	var answer model.Answer
	err := r.db.Preload("Option").Preload("Question.Options").Preload("Submission.Questionnaire").
		Where("submission_id = ? AND question_id = ?", submissionID, questionID).First(&answer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &answer, nil
}

func (r *AnswerRepository) GetAnswersByUserAndQuestionID(questionID model.QuestionID, userID model.UserId) (*[]model.Answer, error) {
	var answers []model.Answer
	result := r.db.Preload("Option").Preload("Question.Options").Preload("Submission.Questionnaire").
		Find(&answers, "question_id = ? AND user_id = ?", questionID, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &answers, nil
}
