package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type NotificationService interface {
	CreateNotification(notification *models.Notification) error
	GetNotificationByID(id uint) (*models.Notification, error)
	GetNotificationsByUser(userID uint) ([]models.Notification, error)
	UpdateNotification(notification *models.Notification) error
	DeleteNotification(id uint) error
}

type notificationService struct {
	repo repositories.NotificationRepository
}

func NewNotificationService(repo repositories.NotificationRepository) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) CreateNotification(notification *models.Notification) error {
	return s.repo.Create(notification)
}

func (s *notificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	return s.repo.FindByID(id)
}

func (s *notificationService) GetNotificationsByUser(userID uint) ([]models.Notification, error) {
	return s.repo.FindByUser(userID)
}

func (s *notificationService) UpdateNotification(notification *models.Notification) error {
	return s.repo.Update(notification)
}

func (s *notificationService) DeleteNotification(id uint) error {
	return s.repo.Delete(id)
}
