package service

import (
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

type IAnswerService interface {
}

type AnswerService struct {
	// dependency injection
	repo domain_repository.IAnswerRepository
}

// NewAnswerService creates a new instance of AnswerService
func NewAnswerService(repo domain_repository.IAnswerRepository) *AnswerService {
	return &AnswerService{repo: repo}
}
