package service

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
)

var (
	ErrNotificationNotFound   = errors.New("notification not found")
	ErrUnauthorizedAccess     = errors.New("unauthorized access to notification")
	ErrNotificationCreateFail = errors.New("failed to create notification")
)

type NotificationService struct {
	repo domain_repository.INotificationRepository
}

// NewNotificationService creates a new NotificationService instance
func NewNotificationService(repo domain_repository.INotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// CreateNotification creates a new notification for a user
func (s *NotificationService) CreateNotification(notification *model.Notification) error {
	return s.repo.CreateNotification(notification)
}

// MarkAsSeen marks a notification as seen if it belongs to the logged-in user
func (s *NotificationService) MarkAsSeen(notificationID, userID uint) error {
	notification, err := s.repo.GetNotificationByID(notificationID)
	if err != nil {
		return ErrNotificationNotFound
	}

	if notification.UserID != userID {
		return ErrUnauthorizedAccess
	}
	return s.repo.SetSeen(notification.ID)
}

// ListNotifications retrieves paginated notifications for a user
func (s *NotificationService) ListNotifications(userID uint, page, limit int) ([]model.Notification, error) {
	return s.repo.GetUserNotifications(userID, page, limit)
}
