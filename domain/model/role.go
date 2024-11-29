package model

// Role represents the roles table
type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;size:255"`
}
