package model

import "database/sql"

// RolePermission represents the role_permission table
type RolePermission struct {
	ID              uint `gorm:"primaryKey"`
	QuestionnaireID uint
	RoleID          uint
	PermissionID    uint
	ExpireAt        sql.NullTime
	Role            Role       `gorm:"foreignKey:RoleID"`
	Permission      Permission `gorm:"foreignKey:PermissionID"`
}
