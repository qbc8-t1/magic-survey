package repository

import (
	"fmt"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *gorm.DB) repository.IUserRepository {
	return &userRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(user *model.User) error {
	fmt.Println(user)
	return r.db.Create(&user).Error
}

// GetUserByID fetches a user by their ID
func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

// DeleteUser deletes a user from the database
func (r *userRepository) DeleteUser(id int) error {
	return r.db.Delete(&model.User{}, id).Error
}