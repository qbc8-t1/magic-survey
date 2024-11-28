package model

// Permission represents the permissions table
type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;size:255"`
	Description string `gorm:"size:500"`
}
