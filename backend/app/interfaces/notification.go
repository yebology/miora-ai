package interfaces

import "miora-ai/app/entities"

// INotificationRepository defines the data access contract for notifications.
type INotificationRepository interface {
	Create(notification *entities.Notification) error
	FindByUser(userID uint) ([]entities.Notification, error)
	MarkAsRead(id, userID uint) error
}

// IEmailClient defines the contract for sending emails.
type IEmailClient interface {
	SendTradeAlert(to, walletAddress, chain, tokenSymbol, direction, value string, liquidity, marketCap float64) error
}
