package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uint) (*models.Task, error)
	FindByProject(projectID uint) ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uint) (*models.Task, error) {
	var t models.Task
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepository) FindByProject(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}
