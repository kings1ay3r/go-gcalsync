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

	t.Run("Successfully inserts calendars and events", func(t *testing.T) {
		anything := mock.Anything
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
			}
			mockToken := &oauth2.Token{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				Expiry:       time.Now().Add(time.Hour),
			}

			mockDao.On("SaveUserTokens", ctx, userID, mockToken.AccessToken, mockToken.RefreshToken, mockToken.Expiry).Return(nil)

			calendarList := &calendar.CalendarList{
				Items: []*calendar.CalendarListEntry{
					{
						Id:      "calendar-id-1",
						Summary: "Test Calendar 1",
					},
				},
			}
			mockCalendarService.On("ListCalendars", ctx).Return(calendarList, nil)

			mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
			mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
			mockTokenSource.On("Token").Return(mockToken, nil)
			mockCalendarService.On("NewService", anything, anything).Return(nil)

			calendars, err := g.FetchCalendars(ctx, userID, code)

			assert.NoError(t, err)
			assert.NotNil(t, calendars)
			assert.Equal(t, 1, len(calendars))
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
			}
			mockconfig.On("Exchange", anything, anything).Return(nil, errors.New("failed to exchange auth code"))

			calendars, err := g.FetchCalendars(ctx, userID, "invalid-auth-code")

			assert.Error(t, err)
			assert.Nil(t, calendars)
			assert.Equal(t, "failed to exchange auth code: failed to exchange auth code", err.Error())

			mockDao.AssertExpectations(t)
			mockconfig.AssertExpectations(t)
		})

		t.Run("should return an error when saving user tokens fails", func(t *testing.T) {
			ctx := context.Background()
			userID := 1
			code := "valid-auth-code"

			mockDao := daomock.NewDAO(t)
			mockCalendarService := mocks.NewCalendarService(t)
			mockconfig := mocks.NewOAuthConfig(t)

			g := &googleCalendar{
				dao:             mockDao,
				calendarService: mockCalendarService,
				config:          mockconfig,
			}
			mockToken := &oauth2.Token{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				Expiry:       time.Now().Add(time.Hour),
			}

			mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
			mockDao.On("SaveUserTokens", ctx, userID, mockToken.AccessToken, mockToken.RefreshToken, mockToken.Expiry).
				Return(errors.New("failed to save user tokens"))

			calendars, err := g.FetchCalendars(ctx, userID, code)

			assert.Error(t, err)
			assert.Nil(t, calendars)
			assert.Equal(t, "failed to save user tokens: failed to save user tokens", err.Error())

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
			}
			mockToken := &oauth2.Token{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				Expiry:       time.Now().Add(time.Hour),
			}

			mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
			mockDao.On("SaveUserTokens", ctx, userID, mockToken.AccessToken, mockToken.RefreshToken, mockToken.Expiry).Return(nil)
			mockconfig.On("TokenSource", anything, mockToken).Return(mockTokenSource, nil)
			mockTokenSource.On("Token").Return(mockToken, nil)
			mockCalendarService.On("NewService", anything, anything).Return(errors.New("failed to create service"))

			calendars, err := g.FetchCalendars(ctx, userID, code)

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
			}
			mockToken := &oauth2.Token{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				Expiry:       time.Now().Add(time.Hour),
			}

			mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
			mockDao.On("SaveUserTokens", ctx, userID, mockToken.AccessToken, mockToken.RefreshToken, mockToken.Expiry).Return(nil)
			mockconfig.On("TokenSource", anything, mockToken).Return(mockTokenSource, nil)
			mockCalendarService.On("NewService", anything, anything).Return(nil)
			mockTokenSource.On("Token").Return(mockToken, nil)
			mockCalendarService.On("ListCalendars", anything).Return(nil, errors.New("failed to fetch calendar list"))

			calendars, err := g.FetchCalendars(ctx, userID, code)

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
			}

			mockToken := &oauth2.Token{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				Expiry:       time.Now().Add(time.Hour),
			}

			mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
			mockDao.On("SaveUserTokens", ctx, userID, mockToken.AccessToken, mockToken.RefreshToken, mockToken.Expiry).Return(errors.New("failed to save tokens"))

			calendars, err := g.FetchCalendars(ctx, userID, code)

			assert.Error(t, err)
			assert.Nil(t, calendars)
			assert.Equal(t, "failed to save user tokens: failed to save tokens", err.Error())

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
			}

			mockToken := &oauth2.Token{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				Expiry:       time.Now().Add(time.Hour),
			}

			mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
			mockDao.On("SaveUserTokens", ctx, userID, mockToken.AccessToken, mockToken.RefreshToken, mockToken.Expiry).Return(nil)
			mockconfig.On("TokenSource", anything, mockToken).Return(mockTokenSource, nil)
			mockTokenSource.On("Token").Return(nil, errors.New("failed to get new token"))

			calendars, err := g.FetchCalendars(ctx, userID, code)

			assert.Error(t, err)
			assert.Nil(t, calendars)
			assert.Equal(t, "failed to get new token", err.Error())

			// Assert that the mocks were called
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
			}

			oldToken := &oauth2.Token{
				AccessToken:  "old-access-token",
				RefreshToken: "old-refresh-token",
				Expiry:       time.Now().Add(-time.Hour), // Expired token
			}

			newToken := &oauth2.Token{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				Expiry:       time.Now().Add(time.Hour),
			}

			calendarList := &calendar.CalendarList{
				Items: []*calendar.CalendarListEntry{
					{
						Id:      "calendar-id-1",
						Summary: "Test Calendar 1",
					},
				},
			}

			mockConfig.On("Exchange", anything, anything).Return(newToken, nil)
			mockDao.On("SaveUserTokens", anything, anything, anything, anything, anything).Return(nil)

			// Simulate fetching the old token
			mockTokenSource.On("Token").Return(oldToken, nil)
			mockConfig.On("TokenSource", anything, newToken).Return(mockTokenSource, nil)
			mockCalendarService.On("NewService", anything, anything).Return(nil)
			mockCalendarService.On("ListCalendars", ctx).Return(calendarList, nil)

			calendars, err := g.FetchCalendars(ctx, userID, code)

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

			g := &googleCalendar{
				dao:             mockDao,
				calendarService: mockCalendarService,
				config:          mockConfig,
			}

			newToken := &oauth2.Token{
				AccessToken: "new-access-token",
				Expiry:      time.Now().Add(time.Hour),
			}

			mockConfig.On("Exchange", anything, anything).Return(newToken, nil)
			// Simulate failure when saving user tokens
			mockDao.On("SaveUserTokens", ctx, userID, newToken.AccessToken, "", newToken.Expiry).Return(errors.New("failed to save user tokens"))

			calendars, err := g.FetchCalendars(ctx, userID, code)

			assert.Error(t, err)
			assert.Nil(t, calendars)
			assert.Equal(t, "failed to save user tokens: failed to save user tokens", err.Error())

			// Assert that the mocks were called
			mockDao.AssertExpectations(t)
			mockConfig.AssertExpectations(t)
		})

	})
}
