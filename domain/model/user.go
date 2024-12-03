package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/QBC8-Team1/magic-survey/pkg/utils"
)

var (
	ErrInvalidUserID = errors.New("userID is required and must be greater than 0")
)

type GenderEnum string
type UserId uint

const (
	Male   GenderEnum = "male"
	Female GenderEnum = "female"
)

// User represents the database model for a user
type User struct {
	ID             uint       `gorm:"primaryKey"`
	FirstName      string     `gorm:"size:255"`
	LastName       string     `gorm:"size:255"`
	Birthdate      string     `gorm:"size:255"`
	City           string     `gorm:"size:255"`
	NationalCode   string     `gorm:"size:10;unique"`
	Gender         GenderEnum `gorm:"type:gender_enum;not null"`
	Email          string     `gorm:"unique;size:255"`
	Password       string     `gorm:"not null"`
	IsActive       bool       `gorm:"not null"`
	Credit         int64
	CreatedAt      time.Time
	UpdatedAt      time.Time       `gorm:"not null"`
	Questionnaires []Questionnaire `gorm:"foreignKey:OwnerID"`
	Notifications  []Notification  `gorm:"foreignKey:UserID"`
	SuperAdmin     *SuperAdmin     `gorm:"foreignKey:UserID"`
}

// TwoFACode stores 2FA codes for users
type TwoFACode struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUserDTO represents the data needed to create a new user
type CreateUserDTO struct {
	FirstName    string     `json:"first_name" validate:"required"`
	LastName     string     `json:"last_name" validate:"required"`
	Email        string     `json:"email" validate:"required,email"`
	NationalCode string     `json:"national_code" validate:"required"`
	Password     string     `json:"password" validate:"required"`
	Gender       GenderEnum `json:"gender" validate:"required,oneof=male female"`
}

// UpdateUserDTO represents the data needed to update an existing user
type UpdateUserDTO struct {
	FirstName    *string     `json:"first_name,omitempty"`
	LastName     *string     `json:"last_name,omitempty"`
	Email        *string     `json:"email,omitempty" validate:"email"`
	NationalCode *string     `json:"national_code,omitempty"`
	Password     *string     `json:"password,omitempty"`
	Gender       *GenderEnum `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
}

// LoginRequest represents user login data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication responses
type AuthResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	TwoFACodeSent bool   `json:"two_fa_code_sent"`
}

// Verify2FACodeRequest validates 2FA code
type Verify2FACodeRequest struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID           UserId `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	NationalCode string `json:"national_code"`
	Gender       string `json:"gender"`
}

// GetFullName returns the full name of a user
func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// ToUserResponse maps a User model to a UserResponse DTO
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:           UserId(user.ID),
		Name:         user.GetFullName(),
		Email:        user.Email,
		NationalCode: user.NationalCode,
		Gender:       string(user.Gender),
	}
}

// ToUserModel maps a CreateUserDTO to a User model
func ToUserModel(dto *CreateUserDTO) *User {
	return &User{
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Email:        dto.Email,
		NationalCode: dto.NationalCode,
		Password:     dto.Password,
		Gender:       dto.Gender,
	}
}

// UpdateUserModel updates the fields of a User model from an UpdateUserDTO
func UpdateUserModel(user *User, dto *UpdateUserDTO) {
	if dto.FirstName != nil {
		user.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		user.LastName = *dto.LastName
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.NationalCode != nil {
		user.NationalCode = *dto.NationalCode
	}
	if dto.Password != nil {
		user.Password = *dto.Password
	}
	if dto.Gender != nil {
		user.Gender = *dto.Gender
	}
}

// Validate checks the User struct for common validation rules.
func (u *User) Validate() error {
	if strings.TrimSpace(u.FirstName) == "" {
		return errors.New("first name is required")
	}
	if strings.TrimSpace(u.LastName) == "" {
		return errors.New("last name is required")
	}
	if !utils.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if len(u.NationalCode) != 10 || !utils.IsAllDigits(u.NationalCode) {
		return errors.New("national code must be a 10-digit number")
	}
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}
