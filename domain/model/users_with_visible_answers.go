package model

import "time"

// UsersWithVisibleAnswers represents the "users_with_visible_answers" table
type UsersWithVisibleAnswers struct {
	ID               uint      `gorm:"primaryKey"`
	RolePermissionID uint      `gorm:"not null;index"`
	UserID           uint      `gorm:"not null;index"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`

	RolePermission RolePermission `gorm:"foreignKey:RolePermissionID"`
	User           User           `gorm:"foreignKey:UserID"`
}
