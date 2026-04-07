package interfaces

import "miora-ai/app/entities"

type INotificationRepository interface {
	Create(n *entities.Notification) error
	FindByUser(userID uint) ([]entities.Notification, error)
	MarkAsRead(id, userID uint) error
}
