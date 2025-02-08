package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	FindByID(id uint) (*models.Review, error)
	FindByProject(projectID uint) ([]models.Review, error)
	Update(review *models.Review) error
	Delete(id uint) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByID(id uint) (*models.Review, error) {
	var review models.Review
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) FindByProject(projectID uint) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Where("project_id = ?", projectID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}
