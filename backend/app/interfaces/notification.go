package interfaces

import "miora-ai/app/entities"

// INotificationRepository defines the data access contract for notifications.
type INotificationRepository interface {
	Create(notification *entities.Notification) error
	FindByUser(userID uint) ([]entities.Notification, error)
	MarkAsRead(id, userID uint) error
}
