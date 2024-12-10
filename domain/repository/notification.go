package domain_repository

import "github.com/QBC8-Team1/magic-survey/domain/model"

// INotificationRepository defines repository methods for notifications
type INotificationRepository interface {
	CreateNotification(notification *model.Notification) error
	GetNotificationByID(id uint) (*model.Notification, error)
	GetUserNotifications(userID uint, limit int, offset int) ([]model.Notification, error)
	SetSeen(notificationID uint) error
	DeleteNotification(id uint) error
}
