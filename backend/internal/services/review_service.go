package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type ReviewService interface {
	CreateReview(review *models.Review) error
	GetReviewByID(id uint) (*models.Review, error)
	GetReviewsByProject(projectID uint) ([]models.Review, error)
	UpdateReview(review *models.Review) error
	DeleteReview(id uint) error
}

type reviewService struct {
	repo repositories.ReviewRepository
}

func NewReviewService(repo repositories.ReviewRepository) ReviewService {
	return &reviewService{repo: repo}
}

func (s *reviewService) CreateReview(review *models.Review) error {
	return s.repo.Create(review)
}

func (s *reviewService) GetReviewByID(id uint) (*models.Review, error) {
	return s.repo.FindByID(id)
}

func (s *reviewService) GetReviewsByProject(projectID uint) ([]models.Review, error) {
	return s.repo.FindByProject(projectID)
}

func (s *reviewService) UpdateReview(review *models.Review) error {
	return s.repo.Update(review)
}

func (s *reviewService) DeleteReview(id uint) error {
	return s.repo.Delete(id)
}
