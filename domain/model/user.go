package model

import (
	"gorm.io/gorm"
)

// User represents the database model for a user
type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;unique"`
	NationalCode string `gorm:"not null;unique"`
	Password     string `gorm:"not null"`
}

// CreateUserDTO represents the data needed to create a new user
type CreateUserDTO struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	NationalCode string `json:"national_code" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

// UpdateUserDTO represents the data needed to update an existing user
type UpdateUserDTO struct {
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty" validate:"email"`
	NationalCode string `json:"national_code,omitempty"`
	Password     string `json:"password,omitempty"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	NationalCode string `json:"national_code"`
}

// ToUserResponse maps a User model to a UserResponse DTO
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		NationalCode: user.NationalCode,
	}
}

// ToUserModel maps a CreateUserDTO to a User model
func ToUserModel(dto *CreateUserDTO) *User {
	return &User{
		Name:         dto.Name,
		Email:        dto.Email,
		NationalCode: dto.NationalCode,
		Password:     dto.Password, // TODO: Hash the password before saving
	}
}

// UpdateUserModel updates the fields of a User model from an UpdateUserDTO
func UpdateUserModel(user *User, dto *UpdateUserDTO) {
	if dto.Name != "" {
		user.Name = dto.Name
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