package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"github.com/QBC8-Team1/magic-survey/pkg/jwt"
	"github.com/QBC8-Team1/magic-survey/pkg/mail"
	t "github.com/QBC8-Team1/magic-survey/pkg/time"
	"github.com/QBC8-Team1/magic-survey/pkg/utils"
	jwt2 "github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidUserID                = errors.New("invalid user ID")
	ErrUserNotFound                 = errors.New("user not found")
	ErrUserOnCreate                 = errors.New("cant Create the user")
	ErrUserOnUpdate                 = errors.New("cant Update the user")
	ErrEmailExists                  = errors.New("mail already exits")
	ErrNationalCodeExists           = errors.New("national code already exits")
	ErrWrongEmailPass               = errors.New("wrong mail or password")
	ErrInvalid2FACode               = errors.New("wrong code")
	ErrCodeExpired                  = errors.New("code expired")
	ErrCodeVerification             = errors.New("cant verify code")
	ErrCantSaveCode                 = errors.New("cant save code")
	ErrCantDeleteCode               = errors.New("cant delete code")
	ErrCantGetCode                  = errors.New("cant get code")
	ErrUserRetrieveFailed           = errors.New("failed to retrieve user")
	ErrUserNotVerified              = errors.New("user is not verified")
	ErrHasOneDayPassedBirthdate     = errors.New("The time to change birthdate has passed and you cannot change it")
	ErrMoreWalletBalanceThanAllowed = errors.New("You are a billionaire, this app will not help you")
	ErrUserIdNotFound               = errors.New("user id not found")
)

type UserService struct {
	repo                        domain_repository.IUserRepository
	authSecret                  string
	expMin, refreshExpMin       uint
	mailPass                    string
	fromMail                    string
	maxSecondForChangeBirthdate int
}

// NewUserService creates a new instance of UserService
func NewUserService(repo domain_repository.IUserRepository, authSecret string, expMin, refreshExpMin uint, mailPass string, fromMail string, maxSecondForChangeBirthdate int) *UserService {
	return &UserService{
		repo:                        repo,
		authSecret:                  authSecret,
		expMin:                      expMin,
		refreshExpMin:               refreshExpMin,
		mailPass:                    mailPass,
		fromMail:                    fromMail,
		maxSecondForChangeBirthdate: maxSecondForChangeBirthdate,
	}
}

// CreateUser handles business logic for creating a user
func (s *UserService) CreateUser(user *model.User) (*model.AuthResponse, error) {
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
	user.IsActive = false

	hashedPassword, err := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	twoFACode := utils.GenerateRandomCode()
	err = mail.SendMail(s.fromMail, s.mailPass, user.Email, "Your 2FA Code", fmt.Sprintf("Your 2FA code is: %s", twoFACode))
	if err != nil {
		return nil, err
	}

	expiry := time.Now().Add(5 * time.Minute)
	err = s.repo.StoreTwoFACode(user.Email, twoFACode, expiry)
	if err != nil {
		return nil, fmt.Errorf("failed to store 2FA code: %w", err)
	}

	return &model.AuthResponse{
		AccessToken:   "",
		RefreshToken:  "",
		TwoFACodeSent: true,
	}, nil
}

// UpdateUser handles business logic for update a user - just [city, first_name, last_name, birthdate]
func (s *UserService) UpdateUser(user *model.User, newUser *model.User) (*model.UserResponse, error) {

	// The user has changed the birthdate - Display error if a day has passed since registration
	if newUser.Birthdate != user.Birthdate && utils.HasTimePassed(user.CreatedAt, s.maxSecondForChangeBirthdate) {
		return nil, ErrHasOneDayPassedBirthdate
	}

	err := s.repo.UpdateUser(newUser)
	if err != nil {
		return nil, err
	}

	return model.ToUserResponse(newUser), nil
}

// IncreaseWalletBalance handles business logic for update a credit user
func (s *UserService) IncreaseWalletBalance(user *model.User, value int64) (*model.UserResponse, error) {
	const maxWalletBalance = 10000000
	newWalletBalance := user.WalletBalance + value
	if newWalletBalance > maxWalletBalance {
		return nil, ErrMoreWalletBalanceThanAllowed
	}
	user.WalletBalance = newWalletBalance
	err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return model.ToUserResponse(user), nil
}

// Show User
//func (s *UserService) ShowUser(id int) (*model.PublicUserResponse, error) {
//	user, err := s.repo.GetUserByID(id)
//	if err != nil {
//		return nil, ErrUserIdNotFound
//	}
//	return model.ToPublicUserResponse(user), nil
//}

// Profile User
func (s *UserService) Profile(user *model.User) (*model.UserResponse, error) {
	return model.ToUserResponse(user), nil
}

// LoginUser handles user logging in logics
func (s *UserService) LoginUser(user *model.LoginRequest) (*model.AuthResponse, error) {
	res, err := s.repo.GetUserByEmail(user.Email)
	if res == nil {
		return nil, ErrWrongEmailPass
	}

	if !res.IsActive {
		return nil, ErrUserNotVerified
	}
	err = utils.CheckPasswordHash(user.Password, res.Password)
	if err != nil {
		return nil, ErrWrongEmailPass
	}

	accessToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(t.AddMinutes(s.expMin, true)),
		},
		UserID: uint(res.ID),
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(t.AddMinutes(s.refreshExpMin, true)),
		},
		UserID: uint(res.ID),
	})

	if err != nil {
		return nil, err
	}
	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}

func (s *UserService) Verify2FACode(userEmail string, enteredCode string) (*model.AuthResponse, error) {
	storedCode, err := s.repo.GetTwoFACode(userEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve 2FA code: %w", err)
	}

	if time.Now().After(storedCode.ExpiresAt) {
		return nil, ErrCodeExpired
	}

	if enteredCode != storedCode.Code {
		return nil, ErrInvalid2FACode
	}
	user, err := s.repo.GetUserByEmail(userEmail)
	accessToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(time.Minute * time.Duration(s.expMin))),
		},
		UserID: uint(user.ID),
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.CreateToken([]byte(s.authSecret), &jwt.UserClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(time.Minute * time.Duration(s.refreshExpMin))),
		},
		UserID: uint(user.ID),
	})
	if err != nil {
		return nil, err
	}

	err = s.repo.RemoveTwoFACode(user.Email)
	if err != nil {
		return nil, ErrCantDeleteCode
	}

	user.IsActive = true
	user.UpdatedAt = time.Now()
	err = s.repo.UpdateUser(user)
	if err != nil {
		return nil, ErrUserOnUpdate
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
