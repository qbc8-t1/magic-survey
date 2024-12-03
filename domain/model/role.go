package model

// Role represents the roles table
type Role struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"unique;size:255"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
