package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type SkillService interface {
	CreateSkill(skill *models.Skill) error
	GetSkillByID(id uint) (*models.Skill, error)
	GetAllSkills() ([]models.Skill, error)
	UpdateSkill(skill *models.Skill) error
	DeleteSkill(id uint) error
}

type skillService struct {
	repo repositories.SkillRepository
}

func NewSkillService(repo repositories.SkillRepository) SkillService {
	return &skillService{repo: repo}
}

func (s *skillService) CreateSkill(skill *models.Skill) error {
	return s.repo.Create(skill)
}

func (s *skillService) GetSkillByID(id uint) (*models.Skill, error) {
	return s.repo.FindByID(id)
}

func (s *skillService) GetAllSkills() ([]models.Skill, error) {
	return s.repo.FindAll()
}

func (s *skillService) UpdateSkill(skill *models.Skill) error {
	return s.repo.Update(skill)
}

func (s *skillService) DeleteSkill(id uint) error {
	return s.repo.Delete(id)
}
