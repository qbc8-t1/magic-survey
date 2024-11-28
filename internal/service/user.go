package service

import (
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

type UserService struct {
	repo domain_repository.IUserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo domain_repository.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser handles business logic for creating a user
func (s *UserService) CreateUser(user *model.User) (*model.User, error) {
	// TODO: handle complex validations here(like querying email is already exists)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := s.repo.CreateUser(user)
	return user, err
}
