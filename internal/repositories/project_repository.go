package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uint) (*models.Project, error)
	FindAll() ([]models.Project, error)
	Update(project *models.Project) error
	Delete(id uint) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) FindAll() ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}
