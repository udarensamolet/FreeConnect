package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"
	"FreeConnect/tests"
)

func TestProjectService(t *testing.T) {
	db := tests.SetupTestDB()
	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)

	// 1) Create a Project
	project := models.Project{
		Title:       "Build an API",
		Description: "Create a REST API in Go",
		Budget:      1500.0,
		Duration:    20,
		ClientID:    1, // If you have a client with ID=1
	}

	err := projectService.CreateProject(&project)
	assert.NoError(t, err)
	assert.NotZero(t, project.ID)

	// 2) Retrieve Project
	got, err := projectService.GetProjectByID(project.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Build an API", got.Title)

	// 3) Update Project
	project.Description = "Updated: Create a REST+GraphQL API"
	err = projectService.UpdateProject(&project)
	assert.NoError(t, err)

	updated, err := projectService.GetProjectByID(project.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated: Create a REST+GraphQL API", updated.Description)

	// 4) Search Projects
	found, err := projectService.SearchProjects("API", "1000", "2000", "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(found), 1)

	// 5) Delete Project
	err = projectService.DeleteProject(project.ID)
	assert.NoError(t, err)

	// 6) Confirm deletion
	_, err = projectService.GetProjectByID(project.ID)
	assert.Error(t, err)
}
