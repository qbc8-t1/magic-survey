package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required" gorm:"unique"`
	NationalCode string `json:"national_code" validate:"required" gorm:"unique"`
	Password     string `json:"password" validate:"required"`
}
