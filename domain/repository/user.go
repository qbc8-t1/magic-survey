package domain_repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"time"
)

// UserRepository is the interface that defines the repository methods
type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByNationalCode(nationalCode string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int) error
	StoreTwoFACode(email string, code string, expiresAt time.Time) error
	GetTwoFACode(email string) (*model.TwoFACode, error)
}
