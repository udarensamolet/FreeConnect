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

func TestNotificationService(t *testing.T) {
	db := tests.SetupTestDB()
	notiRepo := repositories.NewNotificationRepository(db)
	notiService := services.NewNotificationService(notiRepo)

	// 1) Create a Notification
	noti := models.Notification{
		Message:    "Your proposal has been approved!",
		Date:       time.Now(),
		ReadStatus: false,
		Type:       "proposal_update",
		UserID:     1, // if user with ID=1 exists
	}

	err := notiService.CreateNotification(&noti)
	assert.NoError(t, err)
	assert.NotZero(t, noti.ID)

	// 2) Retrieve Notification
	retrieved, err := notiService.GetNotificationByID(noti.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Your proposal has been approved!", retrieved.Message)

	// 3) Update Notification
	noti.ReadStatus = true
	err = notiService.UpdateNotification(&noti)
	assert.NoError(t, err)

	updated, err := notiService.GetNotificationByID(noti.ID)
	assert.NoError(t, err)
	assert.True(t, updated.ReadStatus)

	// 4) Delete Notification
	err = notiService.DeleteNotification(noti.ID)
	assert.NoError(t, err)

	// 5) Confirm deletion
	_, err = notiService.GetNotificationByID(noti.ID)
	assert.Error(t, err)
}
