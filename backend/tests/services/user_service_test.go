package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"
	"FreeConnect/tests"
)

func TestUserService(t *testing.T) {
	db := tests.SetupTestDB()
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	// 1) Register user
	user := models.User{
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  "client",
	}
	err := userService.Register(&user, "somePassword123")
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// 2) Retrieve user
	got, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", got.Name)

	// 3) Update user
	user.Bio = "I am Alice."
	err = userService.UpdateUser(&user)
	assert.NoError(t, err)

	updated, err := userService.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "I am Alice.", updated.Bio)

	// 4) VerifyCredentials
	verifiedUser, err := userService.VerifyCredentials("alice@example.com", "somePassword123")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, verifiedUser.ID)

	// 5) Deletion (optional)
	// There's no "Delete" in userService, but you might do:
	// err = userRepo.Delete(user.ID)
	// assert.NoError(t, err)
}
