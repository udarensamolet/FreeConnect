package services

import (
	"errors"

	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User, plainPassword string) error
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	UpdateUserSkills(userID uint, skillIDs []uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User, plainPassword string) error {
	if existingUser, _ := s.repo.FindByEmail(user.Email); existingUser != nil {
		return errors.New("user with that email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return s.repo.Create(user)
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}

// UpdateUserSkills updates a freelancer's skills (many-to-many relation) automatically.
func (s *userService) UpdateUserSkills(userID uint, skillIDs []uint) error {
	// Retrieve the user.
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// Get the underlying DB from the repository.
	db := s.repo.GetDB()

	// Retrieve skills matching the provided IDs.
	var skills []models.Skill
	if err := db.Where("skill_id IN ?", skillIDs).Find(&skills).Error; err != nil {
		return err
	}

	// Replace the user's skills with the new set.
	return db.Model(user).Association("Skills").Replace(&skills)
}
