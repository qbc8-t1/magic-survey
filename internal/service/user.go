package service

import (
	"fmt"
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
	res, err := s.repo.GetUserByEmail(user.Email)
	fmt.Println("email", res)
	if res != nil {
		return nil, fmt.Errorf("email Already exists")
	}

	res, err = s.repo.GetUserByNationalCode(user.NationalCode)
	fmt.Println("nat", res)

	if res != nil {
		return nil, fmt.Errorf("national code Already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err = s.repo.CreateUser(user)
	return user, err
}
