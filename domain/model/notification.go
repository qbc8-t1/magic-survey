package model

import "time"

// Notification represents the notifications table
type Notification struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Title       string `gorm:"size:255"`
	Description string `gorm:"size:255"`
	IsSeen      bool   `gorm:"default:false"`
	CreatedAt   time.Time
	User        User `gorm:"foreignKey:UserID"`
}
