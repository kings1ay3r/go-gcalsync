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

func TestFetchEvents(t *testing.T) {

	calendarID := "calendar-id"
	userID := 1
	anything := mock.Anything
	// TODO: Refactor other common fields to higher scope / DRY up using method
	// TODO: DRY up assertions

	t.Run("should fetch events successfully", func(t *testing.T) {
		ctx := context.Background()

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

		eventList := &calendar.Events{
			Items: []*calendar.Event{
				{},
			},
		}
		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)
		mockDao.On("SaveUserTokens", ctx, anything, anything, anything, anything, anything).Return(nil)
		mockCalendarService.On("ListEvents", anything, anything, anything).Return(eventList, nil)

		events, err := g.FetchEventsWithCode(ctx, userID, "code", calendarID, "google-account-id")

		assert.NoError(t, err)
		assert.NotNil(t, events)
		assert.Equal(t, 1, len(events))

		mockDao.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})
	t.Run("fetch events fails", func(t *testing.T) {
		ctx := context.Background()

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

		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockDao.On("SaveUserTokens", ctx, anything, anything, anything, anything, anything).Return(nil)
		mockCalendarService.On("ListEvents", anything, anything, anything).Return(nil, errors.New("failed to fetch events"))
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)

		events, err := g.FetchEventsWithCode(ctx, userID, "code", calendarID, "google-account-id")

		assert.Error(t, err)
		assert.Nil(t, events)

		mockDao.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

	t.Run("service init fails", func(t *testing.T) {
		ctx := context.Background()

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

		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockDao.On("SaveUserTokens", ctx, anything, anything, anything, anything, anything).Return(nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, errors.New("failed to init service"))

		events, err := g.FetchEventsWithCode(ctx, userID, "code", calendarID, "google-account-id")

		assert.Error(t, err)
		assert.Nil(t, events)

		mockDao.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})
	t.Run("token refresh fails", func(t *testing.T) {
		ctx := context.Background()

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

		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockDao.On("SaveUserTokens", ctx, anything, anything, anything, anything, anything).Return(nil)
		mockTokenSource.On("Token").Return(nil, errors.New("failed to refresh token"))

		events, err := g.FetchEventsWithCode(ctx, userID, "code", calendarID, "google-account-id")

		assert.Error(t, err)
		assert.Nil(t, events)

		mockDao.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

	t.Run("token save fails", func(t *testing.T) {
		ctx := context.Background()

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

		mockconfig.On("Exchange", anything, anything).Return(mockToken, nil)
		mockDao.On("SaveUserTokens", ctx, anything, anything, anything, anything, anything).Return(errors.New("failed to save token"))

		events, err := g.FetchEventsWithCode(ctx, userID, "code", calendarID, "google-account-id")

		assert.Error(t, err)
		assert.Nil(t, events)

		mockDao.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

	t.Run("token exchange fails", func(t *testing.T) {
		ctx := context.Background()

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
		mockconfig.On("Exchange", anything, anything).Return(nil, errors.New("token exchange failed"))
		//mockDao.On("SaveUserTokens", ctx, anything, anything, anything, anything, anything).Return(errors.New("failed to save token"))

		events, err := g.FetchEventsWithCode(ctx, userID, "code", calendarID, "google-account-id")

		assert.Error(t, err)
		assert.Nil(t, events)

		mockDao.AssertExpectations(t)
		mockCalendarService.AssertExpectations(t)
	})

}
