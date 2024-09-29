package googlecalendar

import (
	"context"
	"errors"
	"gcalsync/gophers/clients/google-calendar/mocks"
	daomock "gcalsync/gophers/dao/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"testing"
	"time"
)

func TestFetchCalendars(t *testing.T) {

	anything := mock.Anything
	dummyjwt := "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiIiLCJhdWQiOiI0NzU1NzY3NTM2Ny1jdmNvNjdxMnQ3Z2RmaHYzcWo4cDE5YmZqam9tcnY1YS5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjExODMiLCJoZCI6Imdvb2dsZS5jb20iLCJlbWFpbCI6ImphbmVAZ29vZ2xlLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhdF9oYXNoIjoiZnNrbGRhYnZmZHMiLCJuYW1lIjoiSmFuZSBEb2UiLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSlRoM0lwcy1Qc0hubVVRa0Fzam8wOWZ1WEtYNUcwT2dtdFVlTjVMSzc4UktYOVBnPXM5Ni1jIiwiZ2l2ZW5fbmFtZSI6IkphbmUiLCJmYW1pbHlfbmFtZSI6IkRvZSIsImlhdCI6MTcyNzYwNzg5OSwiZXhwIjo5NzI3NjExNDk5fQ.GrldkxnJ8ZqO6BpF3F35nGd5VEYDbVTol8XQgpJMhNs"
	mockClientID := "47557675367-cvco67q2t7gdfhv3qj8p19bfjjomrv5a.apps.googleusercontent.com"
	t.Run("should fetch calendars successfully", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		code := "valid-auth-code"

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockconfig := mocks.NewOAuthConfig(t)
		mockTokenSource := mocks.NewTokenSource(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
		}
		mockToken := &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		}

		mockToken = mockToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		calendarList := &calendar.CalendarList{
			Items: []*calendar.CalendarListEntry{
				{
					Id:      "calendar-id-1",
					Summary: "Test Calendar 1",
				},
			},
		}

		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(nil)
		mockCalendarService.On("ListCalendars", anything, anything).Return(calendarList, nil)
		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockconfig.On("GetGoogleAccountID", anything, anything).Return(mockClientID, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)

		calendars, accountID, err := g.FetchCalendars(ctx, userID, code)

		assert.NoError(t, err)
		assert.NotNil(t, calendars)
		assert.Equal(t, 1, len(calendars))
		assert.NotEqual(t, "", accountID)
		assert.Equal(t, "Test Calendar 1", calendars[0].Summary)

		mockDao.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

	t.Run("should return an error on token exchange failure", func(t *testing.T) {
		ctx := context.Background()
		userID := 1

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockconfig := mocks.NewOAuthConfig(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
		}
		mockToken := &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		}

		mockToken = mockToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		mockconfig.On("Exchange", anything, anything).Return(nil, errors.New("failed to exchange auth code"))

		calendars, _, err := g.FetchCalendars(ctx, userID, "invalid-auth-code")

		assert.Error(t, err)
		assert.Nil(t, calendars)
		assert.Equal(t, "failed to exchange auth code", err.Error())

		mockDao.AssertExpectations(t)
		mockconfig.AssertExpectations(t)
	})

	t.Run("should return an error when creating a new calendar service fails", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		code := "valid-auth-code"

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockconfig := mocks.NewOAuthConfig(t)
		mockTokenSource := mocks.NewTokenSource(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
		}
		mockToken := &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		}

		mockToken = mockToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(nil)
		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("TokenSource", anything, mockToken).Return(mockTokenSource, nil)
		mockconfig.On("GetGoogleAccountID", anything, anything).Return(mockClientID, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, errors.New("failed to create service"))

		calendars, _, err := g.FetchCalendars(ctx, userID, code)

		assert.Error(t, err)
		assert.Nil(t, calendars)
		assert.Equal(t, "failed to create service", err.Error())

		mockDao.AssertExpectations(t)
		mockconfig.AssertExpectations(t)
	})

	t.Run("should return an error when fetching the calendar list fails", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		code := "valid-auth-code"

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockconfig := mocks.NewOAuthConfig(t)
		mockTokenSource := mocks.NewTokenSource(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
		}
		mockToken := &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		}

		mockToken = mockToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(nil)
		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("GetGoogleAccountID", anything, anything).Return(mockClientID, nil)
		mockconfig.On("TokenSource", anything, mockToken).Return(mockTokenSource, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)
		mockCalendarService.On("ListCalendars", anything, anything).Return(nil, errors.New("failed to fetch calendar list"))

		calendars, _, err := g.FetchCalendars(ctx, userID, code)

		assert.Error(t, err)
		assert.Nil(t, calendars)
		assert.Equal(t, "failed to fetch calendar list", err.Error())

		mockDao.AssertExpectations(t)
		mockconfig.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

	t.Run("should return an error when saving user tokens fails", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		code := "valid-auth-code"

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockconfig := mocks.NewOAuthConfig(t)
		//mockTokenSource := mocks.NewTokenSource(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
		}

		mockToken := &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		}

		mockToken = mockToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("GetGoogleAccountID", anything, anything).Return(mockClientID, nil)

		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(errors.New("failed to save token"))

		calendars, _, err := g.FetchCalendars(ctx, userID, code)

		assert.Error(t, err)
		assert.Nil(t, calendars)
		assert.Equal(t, "failed to save user tokens: failed to save token", err.Error())

		mockDao.AssertExpectations(t)
		mockconfig.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

	t.Run("should return an error when token source fails to get new token", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		code := "valid-auth-code"

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockconfig := mocks.NewOAuthConfig(t)
		mockTokenSource := mocks.NewTokenSource(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
		}
		mockToken := &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		}

		mockToken = mockToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(nil).Once()
		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockTokenSource.On("Token").Return(nil, errors.New("failed to retrieve user tokens"))
		mockconfig.On("GetGoogleAccountID", anything, anything).Return(mockClientID, nil)

		calendars, _, err := g.FetchCalendars(ctx, userID, code)

		assert.Error(t, err)
		assert.Nil(t, calendars)
		assert.Equal(t, "failed to refresh tokens: failed to retrieve user tokens", err.Error())

		mockDao.AssertExpectations(t)
		mockconfig.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

	t.Run("should correctly save and refresh tokens when they change", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		code := "valid-auth-code"

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockConfig := mocks.NewOAuthConfig(t)
		mockTokenSource := mocks.NewTokenSource(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockConfig,
			clientCache:     NewClientCache(0),
		}

		oldToken := &oauth2.Token{
			AccessToken:  "old-access-token",
			RefreshToken: "old-refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		}

		oldToken = oldToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		newToken := &oauth2.Token{
			AccessToken:  "new-access-token2",
			RefreshToken: "new-refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		}

		calendarList := &calendar.CalendarList{
			Items: []*calendar.CalendarListEntry{
				{
					Id:      "calendar-id-1",
					Summary: "Test Calendar 1",
				},
			},
		}

		mockConfig.On("Exchange", anything, anything).Return(oldToken, nil)
		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(nil)

		mockTokenSource.On("Token").Return(newToken, nil)
		mockConfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockConfig.On("GetGoogleAccountID", anything, anything).Return(mockClientID, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)
		mockCalendarService.On("ListCalendars", ctx, anything).Return(calendarList, nil)

		calendars, _, err := g.FetchCalendars(ctx, userID, code)

		assert.NoError(t, err)
		assert.NotNil(t, calendars)

		// Assert that the updated token is saved
		mockDao.AssertExpectations(t)
		mockConfig.AssertExpectations(t)
		mockTokenSource.AssertExpectations(t)
	})

	t.Run("should return an error when saving user tokens fails during token refresh", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		code := "valid-auth-code"

		mockDao := daomock.NewDAO(t)
		mockCalendarService := mocks.NewCalendarService(t)
		mockConfig := mocks.NewOAuthConfig(t)
		mockTokenSource := mocks.NewTokenSource(t)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockConfig,
			clientCache:     NewClientCache(0),
		}

		oldToken := &oauth2.Token{
			AccessToken:  "old-access-token",
			RefreshToken: "old-refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		}

		oldToken = oldToken.WithExtra(map[string]interface{}{
			"id_token": dummyjwt,
		})

		newToken := &oauth2.Token{
			AccessToken:  "new-access-token2",
			RefreshToken: "new-refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		}

		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(nil).Once()
		mockConfig.On("Exchange", anything, anything).Return(oldToken, nil)
		mockTokenSource.On("Token").Return(newToken, nil)
		mockConfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockConfig.On("GetGoogleAccountID", anything, anything).Return(mockClientID, nil)
		mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything, anything).Return(errors.New("failed to save user tokens"))

		calendars, _, err := g.FetchCalendars(ctx, userID, code)

		assert.Error(t, err)
		assert.Nil(t, calendars)
		assert.Equal(t, "failed to save refreshed token: failed to save user tokens", err.Error())

		// Assert that the mocks were called
		mockDao.AssertExpectations(t)
		mockConfig.AssertExpectations(t)
	})

}
