package model

// RoleUser represents the role_user table
type RoleUser struct {
	RoleID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
	Role   Role `gorm:"foreignKey:RoleID"`
	User   User `gorm:"foreignKey:UserID"`
}
