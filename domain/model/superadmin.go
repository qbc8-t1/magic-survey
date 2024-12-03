package model

import "time"

// SuperAdmin represents the super_admins table
type Superadmin struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	GrantedBy *uint
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}
