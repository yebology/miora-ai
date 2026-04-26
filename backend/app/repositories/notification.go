package repositories

import (
	"miora-ai/app/entities"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {

	return &NotificationRepository{db: db}

}

func (r *NotificationRepository) Create(n *entities.Notification) error {

	return r.db.Create(n).Error

}

func (r *NotificationRepository) FindByUser(userID uint) ([]entities.Notification, error) {

	var items []entities.Notification
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil

}

func (r *NotificationRepository) MarkAsRead(id, userID uint) error {

	return r.db.Model(&entities.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("read", true).Error

}
