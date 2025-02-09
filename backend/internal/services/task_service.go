package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type TaskService interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id uint) (*models.Task, error)
	GetTasksByProject(projectID uint) ([]models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id uint) error
}

type taskService struct {
	repo repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(task *models.Task) error {
	return s.repo.Create(task)
}

func (s *taskService) GetTaskByID(id uint) (*models.Task, error) {
	return s.repo.FindByID(id)
}

func (s *taskService) GetTasksByProject(projectID uint) ([]models.Task, error) {
	return s.repo.FindByProject(projectID)
}

func (s *taskService) UpdateTask(task *models.Task) error {
	return s.repo.Update(task)
}

func (s *taskService) DeleteTask(id uint) error {
	return s.repo.Delete(id)
}
