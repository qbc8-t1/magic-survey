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

// CreateNotificationDTO represents the data needed to create a new notification
type CreateNotificationDTO struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// NotificationResponseDTO represents the notification data returned in API responses
type NotificationResponseDTO struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsSeen      bool   `json:"is_seen"`
	CreatedAt   string `json:"created_at"`
}

func ToNotificationResponseDTO(notification *Notification) *NotificationResponseDTO {
	return &NotificationResponseDTO{
		ID:          notification.ID,
		UserID:      notification.UserID,
		Title:       notification.Title,
		Description: notification.Description,
		IsSeen:      notification.IsSeen,
		CreatedAt:   notification.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToNotificationModel(dto *CreateNotificationDTO, UserID uint) *Notification {
	return &Notification{
		UserID:      UserID,
		Title:       dto.Title,
		Description: dto.Description,
		IsSeen:      false,
	}
}
