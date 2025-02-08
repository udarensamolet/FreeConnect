package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type SkillRepository interface {
	Create(skill *models.Skill) error
	FindByID(id uint) (*models.Skill, error)
	FindAll() ([]models.Skill, error)
	Update(skill *models.Skill) error
	Delete(id uint) error
}

type skillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepository{db: db}
}

func (r *skillRepository) Create(skill *models.Skill) error {
	return r.db.Create(skill).Error
}

func (r *skillRepository) FindByID(id uint) (*models.Skill, error) {
	var skill models.Skill
	if err := r.db.First(&skill, id).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) FindAll() ([]models.Skill, error) {
	var skills []models.Skill
	if err := r.db.Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *skillRepository) Update(skill *models.Skill) error {
	return r.db.Save(skill).Error
}

func (r *skillRepository) Delete(id uint) error {
	return r.db.Delete(&models.Skill{}, id).Error
}
