package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"gorm.io/gorm"
	"time"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of userRepository
func NewUserRepository(db *gorm.DB) domain_repository.IUserRepository {
	return &userRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(&user).Error
}

// GetUserByID fetches a user by their ID
func (r *userRepository) GetUserByID(id model.UserId) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail fetches a user by their Email
func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByNationalCode fetches a user by their Email
func (r *userRepository) GetUserByNationalCode(nationalCode string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("national_code = ?", nationalCode).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

// DeleteUser deletes a user from the database
func (r *userRepository) DeleteUser(id model.UserId) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) StoreTwoFACode(email string, code string, expiresAt time.Time) error {
	twoFACode := &model.TwoFACode{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(twoFACode).Error; err != nil {
		return service.ErrCantSaveCode
	}

	return nil
}

func (r *userRepository) GetTwoFACode(email string) (*model.TwoFACode, error) {
	twoFACode := &model.TwoFACode{}
	if err := r.db.Where("email = ?", email).Order("created_at desc").First(&twoFACode).Error; err != nil {
		return nil, service.ErrCantGetCode
	}

	return twoFACode, nil
}

func (r *userRepository) RemoveTwoFACode(email string) error {
	if err := r.db.Where("email = ?", email).Delete(&model.TwoFACode{}).Error; err != nil {
		return service.ErrCantDeleteCode
	}

	return nil
}
