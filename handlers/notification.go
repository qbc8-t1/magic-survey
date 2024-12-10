package handlers

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/QBC8-Team1/magic-survey/pkg/response"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func CreateNotification(notificationService service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse request body
		var dto model.CreateNotificationDTO
		if err := c.BodyParser(&dto); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid request body", err)
		}
		user := c.Locals("user").(model.User)
		notification := model.ToNotificationModel(&dto, user.ID)

		if err := notificationService.CreateNotification(notification); err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Failed to create notification", err)
		}

		return response.Success(c, fiber.StatusCreated, "Notification created successfully", nil)
	}
}

func MarkNotificationAsSeen(notificationService service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get notification ID from URL params
		notificationID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid notification ID", err)
		}

		userID := c.Locals("user").(model.User).ID // Logged-in user's ID

		if err := notificationService.MarkAsSeen(uint(notificationID), userID); err != nil {
			return response.Error(c, fiber.StatusForbidden, err.Error(), nil)
		}

		return response.Success(c, fiber.StatusOK, "Notification marked as seen", nil)
	}
}

/*func ListNotifications(notificationService service.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse pagination parameters
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		userID := c.Locals("user").(model.User).ID // Logged-in user's ID

		notifications, err := notificationService.ListNotifications(userID, page, limit)
		if err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch notifications", err)
		}

		return response.Success(c, fiber.StatusOK, "Notifications fetched successfully", notifications)
	}
}*/
