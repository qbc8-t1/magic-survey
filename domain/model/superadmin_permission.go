package model

type SuperadminPermission struct {
	ID           uint `gorm:"primaryKey"`
	SuperadminID uint `gorm:"foreignKey:SuperadminID"`
	PermissionID uint `gorm:"foreignKey:PermissionID"`
	Superadmin   Superadmin
	Permission   Permission
}
