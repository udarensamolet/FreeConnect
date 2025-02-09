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
	SearchProjects(search, minBudgetStr, maxBudgetStr, status string) ([]models.Project, error)
}

type projectRepository struct {
	db *gorm.DB
}

func (r *projectRepository) SearchProjects(search, minBudgetStr, maxBudgetStr, status string) ([]models.Project, error) {
	db := r.db

	if search != "" {
		// Filter by title or description (LIKE):
		db = db.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if minBudgetStr != "" {
		// Convert to float
		db = db.Where("budget >= ?", minBudgetStr)
	}

	if maxBudgetStr != "" {
		// Convert to float
		db = db.Where("budget <= ?", maxBudgetStr)
	}

	if status != "" {
		// if you want a single status or maybe multiple
		db = db.Where("status = ?", status)
	}

	var projects []models.Project
	if err := db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	// If you want to preload tasks
	if err := r.db.Preload("Tasks").First(&project, id).Error; err != nil {
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
