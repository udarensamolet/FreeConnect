package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type ProjectService interface {
	CreateProject(project *models.Project) error
	GetProjectByID(id uint) (*models.Project, error)
	GetAllProjects() ([]models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(id uint) error
}

type projectService struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(project *models.Project) error {
	// Optionally, validate project fields (e.g., budget > 0, duration > 0)
	return s.repo.Create(project)
}

func (s *projectService) GetProjectByID(id uint) (*models.Project, error) {
	return s.repo.FindByID(id)
}

func (s *projectService) GetAllProjects() ([]models.Project, error) {
	return s.repo.FindAll()
}

func (s *projectService) UpdateProject(project *models.Project) error {
	return s.repo.Update(project)
}

func (s *projectService) DeleteProject(id uint) error {
	return s.repo.Delete(id)
}
