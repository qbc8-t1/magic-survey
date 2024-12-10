package repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new instance of notificationRepository
func NewNotificationRepository(db *gorm.DB) domain_repository.INotificationRepository {
	return &notificationRepository{db: db}
}

// CreateNotification creates a new notification in the database
func (r *notificationRepository) CreateNotification(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

// GetNotificationByID fetches a notification by its ID
func (r *notificationRepository) GetNotificationByID(id uint) (*model.Notification, error) {
	var notification model.Notification
	if err := r.db.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetUserNotifications fetches all notifications for a specific user
func (r *notificationRepository) GetUserNotifications(userID uint, limit int, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// SetSeen marks a notification as seen
func (r *notificationRepository) SetSeen(notificationID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ?", notificationID).
		Update("is_seen", true).Error
}

// DeleteNotification deletes a notification from the database
func (r *notificationRepository) DeleteNotification(id uint) error {
	return r.db.Delete(&model.Notification{}, id).Error
}
