package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/QBC8-Team1/magic-survey/pkg/utils"
)

type GenderEnum string
type UserId uint

const (
	Male   GenderEnum = "male"
	Female GenderEnum = "female"
)

// User represents the database model for a user
type User struct {
	ID             uint        `gorm:"primaryKey"`
	FirstName      string      `gorm:"size:255"`
	LastName       string      `gorm:"size:255"`
	Birthdate      string      `gorm:"size:255"`
	City           string      `gorm:"size:255"`
	NationalCode   string      `gorm:"size:10;unique"`
	Gender         *GenderEnum `gorm:"type:gender_enum"`
	Email          string      `gorm:"unique;size:255"`
	Password       string      `gorm:"not null"`
	Credit         int64
	CreatedAt      time.Time
	UpdatedAt      time.Time       `gorm:"not null"`
	Questionnaires []Questionnaire `gorm:"foreignKey:OwnerID"`
	Notifications  []Notification  `gorm:"foreignKey:UserID"`
	SuperAdmin     *SuperAdmin     `gorm:"foreignKey:UserID"`
}

func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// CreateUserDTO represents the data needed to create a new user
type CreateUserDTO struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	NationalCode string `json:"national_code" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

// UpdateUserDTO represents the data needed to update an existing user
type UpdateUserDTO struct {
	FirstName    string `json:"last_name,omitempty"`
	LastName     string `json:"first_name,omitempty"`
	Email        string `json:"email,omitempty" validate:"email"`
	NationalCode string `json:"national_code,omitempty"`
	Password     string `json:"password,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID           UserId `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	NationalCode string `json:"national_code"`
}

// ToUserResponse maps a User model to a UserResponse DTO
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:           UserId(user.ID),
		Name:         user.GetFullName(),
		Email:        user.Email,
		NationalCode: user.NationalCode,
	}
}

// ToUserModel maps a CreateUserDTO to a User model
func ToUserModel(dto *CreateUserDTO) *User {
	return &User{
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Email:        dto.Email,
		NationalCode: dto.NationalCode,
		Password:     dto.Password, // TODO: Hash the password before saving
	}
}

// UpdateUserModel updates the fields of a User model from an UpdateUserDTO
func UpdateUserModel(user *User, dto *UpdateUserDTO) {
	if dto.FirstName != "" {
		user.FirstName = dto.FirstName
	}
	if dto.LastName != "" {
		user.LastName = dto.LastName
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.NationalCode != "" {
		user.NationalCode = dto.NationalCode
	}
	if dto.Password != "" {
		user.Password = dto.Password // TODOs: Hash the password before saving
	}
}

// Validate checks the User struct for common validation rules.
func (u *User) Validate() error {
	// Validate FirstName
	if strings.TrimSpace(u.FirstName) == "" {
		return errors.New("first name is required")
	}

	// Validate LastName
	if strings.TrimSpace(u.LastName) == "" {
		return errors.New("last name is required")
	}

	// Validate Email
	if !utils.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	// Validate NationalCode
	if len(u.NationalCode) != 10 || !utils.IsAllDigits(u.NationalCode) {
		return errors.New("national code must be a 10-digit number")
	}

	// Validate Password
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}
