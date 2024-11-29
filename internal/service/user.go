package service

import (
	"fmt"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"github.com/QBC8-Team1/magic-survey/pkg/utils"
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
	res, err := s.repo.GetUserByEmail(user.Email)
	if res != nil {
		return nil, fmt.Errorf("email Already exists")
	}

	res, err = s.repo.GetUserByNationalCode(user.NationalCode)

	if res != nil {
		return nil, fmt.Errorf("national code Already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err = s.repo.CreateUser(user)
	return user, err
}

// LoginUser handles user logging in logics
func (s *UserService) LoginUser(user *model.LoginRequest) (*model.User, error) {
	res, err := s.repo.GetUserByEmail(user.Email)
	fmt.Println(user.Password, res.Password)
	err = utils.CheckPasswordHash(user.Password, res.Password)
	if res != nil {
		return nil, fmt.Errorf("email or password is wrong!")
	}
	return res, err
}
