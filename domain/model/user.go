package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/QBC8-Team1/magic-survey/pkg/utils"
)

var (
	ErrInvalidUserIDCreate = errors.New("userID is required and must be greater than 0")
)

type GenderEnum string
type UserID uint

const (
	Male   GenderEnum = "male"
	Female GenderEnum = "female"
)

// User represents the database model for a user
type User struct {
	ID                     uint        `gorm:"primaryKey"`
	FirstName              string      `gorm:"size:255"`
	LastName               string      `gorm:"size:255"`
	Birthdate              string      `gorm:"size:255"`
	City                   string      `gorm:"size:255"`
	NationalCode           string      `gorm:"size:10;unique"`
	Gender                 *GenderEnum `gorm:"type:gender_enum"`
	Email                  string      `gorm:"unique;size:255"`
	Password               string      `gorm:"not null"`
	IsActive               bool        `gorm:"not null"`
	WalletBalance          int64
	MaxQuestionnairesCount int `gorm:"null"`
	CreatedAt              time.Time
	UpdatedAt              time.Time       `gorm:"not null"`
	Questionnaires         []Questionnaire `gorm:"foreignKey:OwnerID"`
	Notifications          []Notification  `gorm:"foreignKey:UserID"`
	Superadmin             *Superadmin     `gorm:"foreignKey:UserID"`
	Roles                  []Role          `gorm:"many2many:role_users;"`
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
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	NationalCode string `json:"national_code" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

// UpdateUserDTO represents the data needed to update an existing user
type UpdateUserDTO struct {
	FirstName string     `json:"first_name" validate:"required"`
	LastName  string     `json:"last_name" validate:"required"`
	Birthdate string     `json:"birthdate" validate:"required"`
	Gender    GenderEnum `json:"gender" validate:"required,oneof=male female"`
	City      string     `json:"city" validate:"required"`
}

// IncreaseWalletBalanceDTO represents the data needed to update credit user
type IncreaseWalletBalanceDTO struct {
	Value string `json:"value" validate:"required"`
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
	ID            UserID `json:"id"`
	Name          string `json:"name"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	NationalCode  string `json:"national_code"`
	Gender        string `json:"gender"`
	Birthdate     string `json:"birthdate"`
	City          string `json:"city"`
	WalletBalance int64  `json:"wallet_balance"`
}

// PublicUserResponse represents the user data returned in API responses
type PublicUserResponse struct {
	ID     UserID `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

// GetFullName returns the full name of a user
func (u *User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) GetGender() string {
	if u.Gender == nil {
		return "unknown"
	}
	return string(*u.Gender)
}

// ToUserResponse maps a User model to a UserResponse DTO
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:            UserID(user.ID),
		Name:          user.GetFullName(),
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		NationalCode:  user.NationalCode,
		Gender:        user.GetGender(),
		City:          user.City,
		Birthdate:     user.Birthdate,
		WalletBalance: user.WalletBalance,
	}
}

// ToPublicUserResponse maps a User model to a UserResponse DTO
func ToPublicUserResponse(user *User) *PublicUserResponse {
	return &PublicUserResponse{
		ID:     UserID(user.ID),
		Name:   user.GetFullName(),
		Gender: user.GetGender(),
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
	}
}

// ToUserModelForUpdate maps a UpdateUserDTO to a User model
func ToUserModelForUpdate(user User, dto *UpdateUserDTO) User {
	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.Birthdate = dto.Birthdate
	user.City = dto.City
	return user
}

// UpdateUserModel updates the fields of a User model from an UpdateUserDTO
/*func UpdateUserModel(user *User, dto *UpdateUserDTO) {
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
}*/

//}

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

	isValidNationalCode, err := utils.IsValidNationalCode(u.NationalCode)
	if err != nil {
		return errors.New("national code validation failed")
	}
	if !isValidNationalCode {
		return errors.New("national code is not valid")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	// isValid, message := utils.IsValidBirthdate(u.Birthdate)
	// if !isValid {
	// 	return errors.New("birthdate - " + message)
	// }

	// isValid, message = utils.IsValidCity(u.City)
	// if !isValid {
	// 	return errors.New("city - " + message)
	// }

	return nil
}
