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
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(user *models.User, plainPassword string) error {
	// Check for duplicate email.
	if existingUser, _ := s.repo.FindByEmail(user.Email); existingUser != nil {
		return errors.New("user with that email already exists")
	}

	// Hash the plain-text password.
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
