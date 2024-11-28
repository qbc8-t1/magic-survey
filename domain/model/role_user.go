package model

// RoleUser represents the role_user table
type RoleUser struct {
	ID     uint `gorm:"primaryKey"`
	RoleID uint
	UserID uint
	Role   Role `gorm:"foreignKey:RoleID"`
	User   User `gorm:"foreignKey:UserID"`
}
