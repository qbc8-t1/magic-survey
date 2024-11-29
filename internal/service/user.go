package service

import (
	"errors"
	"github.com/QBC8-Team1/magic-survey/pkg/jwt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"github.com/QBC8-Team1/magic-survey/pkg/utils"
)

var (
	ErrUserOnCreate       = errors.New("Cant Create the user")
	ErrEmailExists        = errors.New("email already exits")
	ErrNationalCodeExists = errors.New("national code already exits")
	ErrWrongEmailPass     = errors.New("wrong email or password")
)

type UserService struct {
	repo                  domain_repository.IUserRepository
	authSecret            string
	expMin, refreshExpMin uint
}

// NewUserService creates a new instance of UserService
func NewUserService(repo domain_repository.IUserRepository, authSecret string, expMin, refreshExpMin uint) *UserService {
	return &UserService{
		repo:          repo,
		authSecret:    authSecret,
		expMin:        expMin,
		refreshExpMin: refreshExpMin}
}

// CreateUser handles business logic for creating a user
func (s *UserService) CreateUser(user *model.User) (*model.User, error) {
	res, err := s.repo.GetUserByEmail(user.Email)
	if res != nil {
		return nil, ErrEmailExists
	}

	res, err = s.repo.GetUserByNationalCode(user.NationalCode)

	if res != nil {
		return nil, ErrNationalCodeExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashedPassword, err := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	err = s.repo.CreateUser(user)

	accessToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.AddMinutes(s.expMin, true)),
		},
		UserID: uint(userID),
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.AddMinutes(s.refreshExpMin, true)),
		},
		UserID: uint(user.ID),
	})

	if err != nil {
		return nil, err
	}
	return user, err
}

// LoginUser handles user logging in logics
func (s *UserService) LoginUser(user *model.LoginRequest) (*model.User, error) {
	res, err := s.repo.GetUserByEmail(user.Email)
	if res == nil {
		return nil, ErrWrongEmailPass
	}

	err = utils.CheckPasswordHash(user.Password, res.Password)
	if err != nil {
		return nil, ErrWrongEmailPass
	}
	return res, err
}
