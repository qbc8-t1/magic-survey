package repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

// UserRepository is the interface that defines the repository methods
type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id int) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int) error
}
