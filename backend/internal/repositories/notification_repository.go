package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByID(id uint) (*models.Notification, error)
	FindByUser(userID uint) ([]models.Notification, error)
	Update(notification *models.Notification) error
	Delete(id uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByID(id uint) (*models.Notification, error) {
	var n models.Notification
	if err := r.db.First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) FindByUser(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	if err := r.db.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Notification{}, id).Error
}
