package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"
	"FreeConnect/tests"
)

func TestSkillService(t *testing.T) {
	db := tests.SetupTestDB()
	skillRepo := repositories.NewSkillRepository(db)
	skillService := services.NewSkillService(skillRepo)

	// 1) Create a Skill
	skill := models.Skill{
		Name:        "Go Programming",
		Level:       "Intermediate",
		Description: "Golang programming language expertise",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := skillService.CreateSkill(&skill)
	assert.NoError(t, err)
	assert.NotZero(t, skill.ID)

	// 2) Retrieve
	fetched, err := skillService.GetSkillByID(skill.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Go Programming", fetched.Name)

	// 3) Update
	skill.Description = "Updated: advanced concurrency skills"
	err = skillService.UpdateSkill(&skill)
	assert.NoError(t, err)

	updated, err := skillService.GetSkillByID(skill.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated: advanced concurrency skills", updated.Description)

	// 4) Delete
	err = skillService.DeleteSkill(skill.ID)
	assert.NoError(t, err)

	// 5) Confirm
	_, err = skillService.GetSkillByID(skill.ID)
	assert.Error(t, err)
}
