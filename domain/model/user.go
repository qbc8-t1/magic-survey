package model

type User struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required"`
	NationalCode string `json:"national_code" validate:"required"`
	Password     string `json:"password" validate:"required"`
}
