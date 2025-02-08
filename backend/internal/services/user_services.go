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

	// NEW: For authentication
	VerifyCredentials(email, password string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register registers a user (client, freelancer, or admin)
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

// GetUserByID retrieves a user by its ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}

// UpdateUserSkills updates a freelancer's skill set
func (s *userService) UpdateUserSkills(userID uint, skillIDs []uint) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	db := s.repo.GetDB()

	var skills []models.Skill
	if err := db.Where("skill_id IN ?", skillIDs).Find(&skills).Error; err != nil {
		return err
	}
	return db.Model(user).Association("Skills").Replace(&skills)
}

// VerifyCredentials checks user email+password, returning an error if invalid
func (s *userService) VerifyCredentials(email, password string) (*models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	// Compare hashed password
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
