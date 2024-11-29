package model

import "time"

// RolePermission represents the role_permission table
type RolePermission struct {
	ID              uint `gorm:"primaryKey"`
	QuestionnaireID uint
	RoleID          uint
	PermissionID    uint
	ExpireAt        *time.Time
	Role            Role       `gorm:"foreignKey:RoleID"`
	Permission      Permission `gorm:"foreignKey:PermissionID"`
}
