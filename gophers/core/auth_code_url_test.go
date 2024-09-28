package core

import (
	"context"
	"fmt"
	"gcalsync/gophers/clients/google-calendar/mocks"
	"gcalsync/gophers/middlewares/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAuthCodeURL(t *testing.T) {
	t.Run("Successfully returns auth code URL", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), auth.ContextUserIDKey, "mockUserID")
		mockGoogleCalClient := mocks.NewGoogleCalendar(t)
		calendarClient := calendarClient{
			googleCalClient: mockGoogleCalClient,
		}

		expectedAuthURL := "https://example.com/auth"

		mockGoogleCalClient.On("GetAuthCodeURL", ctx, "mockUserID").Return(expectedAuthURL)

		authURL, err := calendarClient.GetAuthCodeURL(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedAuthURL, authURL)
		mockGoogleCalClient.AssertCalled(t, "GetAuthCodeURL", ctx, "mockUserID")
	})

	t.Run("Error when user ID is not present in context", func(t *testing.T) {
		ctx := context.Background()
		mockGoogleCalClient := mocks.NewGoogleCalendar(t)
		calendarClient := calendarClient{
			googleCalClient: mockGoogleCalClient,
		}

		authURL, err := calendarClient.GetAuthCodeURL(ctx)

		assert.Error(t, err)
		assert.Equal(t, "", authURL)
		assert.Equal(t, fmt.Errorf("failed to get auth code url").Error(), err.Error())
		mockGoogleCalClient.AssertNotCalled(t, "GetAuthCodeURL", mock.Anything, mock.Anything)
	})
}
