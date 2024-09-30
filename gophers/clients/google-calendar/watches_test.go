package googlecalendar

import (
	"context"
	"errors"
	"gcalsync/gophers/clients/google-calendar/mocks"
	mocks2 "gcalsync/gophers/dao/mocks"
	"gcalsync/gophers/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"testing"
	"time"
)

const anything = mock.Anything

var someErr = errors.New("some error")

func Test_RenewWatches(t *testing.T) {

	expirationTime := time.Now().Add(__GOOGLE_WATCH_LIFE__).UnixMilli()
	mockChannel := &calendar.Channel{
		Expiration: expirationTime,
		ResourceId: "resource_1",
	}
	watch := &dto.WatchData{
		ID:               1,
		UserID:           123,
		CalendarID:       456,
		GoogleCalendarID: "primary",
		AccountID:        "account_1",
		Expiry:           time.Now().Add(-time.Hour), // Watch has expired
	}
	mockGAPIError := &googleapi.Error{
		Code:    0,
		Message: "",
		Details: []interface{}{
			map[string]interface{}{
				"Reason": __GOOGLE_UNSUPPORTED_CALENDAR_FOR_WATCH_REASON__,
			},
		},
		Body:   "",
		Header: nil,
		Errors: nil,
	}

	t.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		mockCalendarService := mocks.NewCalendarService(t)
		mockDao := mocks2.NewDAO(t)
		mockTokenSource := mocks.NewTokenSource(t)
		mockconfig := mocks.NewOAuthConfig(t)

		watches := []*dto.WatchData{
			watch,
		}

		token := &oauth2.Token{AccessToken: "test-access-token"}

		tokens := map[string]*oauth2.Token{
			"account_1": token,
		}

		mockCalendarService.On("Watch", anything, anything, anything, anything).Return(mockChannel, nil)
		mockDao.On("GetUserTokens", anything, anything, anything, anything).Return(tokens["account_1"], nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockTokenSource.On("Token").Return(token, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
			webhookURL:      "example.com",
		}

		res, errs := g.RenewWatches(ctx, watches, tokens)

		assert.Empty(t, errs, "Expected no errors")
		assert.Len(t, res.Succesful, 1, "Expected one watch to be successfully renewed")
		assert.Equal(t, watches[0].ID, res.Succesful[0].ID, "Watch IDs should match")
		assert.Equal(t, time.UnixMilli(expirationTime), res.Succesful[0].Expiry, "Expiry should match expected time")

	})

	t.Run("watch initialization fails", func(t *testing.T) {
		ctx := context.Background()

		mockCalendarService := mocks.NewCalendarService(t)
		mockDao := mocks2.NewDAO(t)
		mockTokenSource := mocks.NewTokenSource(t)
		mockconfig := mocks.NewOAuthConfig(t)

		watches := []*dto.WatchData{
			watch,
		}

		token := &oauth2.Token{AccessToken: "test-access-token"}

		tokens := map[string]*oauth2.Token{
			"account_1": token,
		}

		mockCalendarService.On("Watch", anything, anything, anything, anything).Return(nil, someErr)
		mockDao.On("GetUserTokens", anything, anything, anything, anything).Return(tokens["account_1"], nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockTokenSource.On("Token").Return(token, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
			webhookURL:      "example.com",
		}

		res, errs := g.RenewWatches(ctx, watches, tokens)

		assert.NotEmpty(t, errs, "Expected  errors")
		assert.Len(t, res.Failed, 1, "Expected one watch to fail")
		assert.Equal(t, watches[0].ID, res.Failed[0].ID, "Watch IDs should match")

	})

	t.Run("watch initialization fails for static calendar", func(t *testing.T) {
		ctx := context.Background()

		mockCalendarService := mocks.NewCalendarService(t)
		mockDao := mocks2.NewDAO(t)
		mockTokenSource := mocks.NewTokenSource(t)
		mockconfig := mocks.NewOAuthConfig(t)

		watches := []*dto.WatchData{
			watch,
		}

		token := &oauth2.Token{AccessToken: "test-access-token"}

		tokens := map[string]*oauth2.Token{
			"account_1": token,
		}

		mockCalendarService.On("Watch", anything, anything, anything, anything).Return(nil, mockGAPIError)
		mockDao.On("GetUserTokens", anything, anything, anything, anything).Return(tokens["account_1"], nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockTokenSource.On("Token").Return(token, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
			webhookURL:      "example.com",
		}

		res, errs := g.RenewWatches(ctx, watches, tokens)

		assert.NotEmpty(t, errs, "Expected  errors")
		assert.Len(t, res.Failed, 1, "Expected one watch to fail")
		assert.Equal(t, watches[0].ID, res.Failed[0].ID, "Watch IDs should match")

	})

	t.Run("watch initialization fails partially", func(t *testing.T) {
		ctx := context.Background()

		mockCalendarService := mocks.NewCalendarService(t)
		mockDao := mocks2.NewDAO(t)
		mockTokenSource := mocks.NewTokenSource(t)
		mockconfig := mocks.NewOAuthConfig(t)

		watches := []*dto.WatchData{
			watch,
			watch,
		}

		token := &oauth2.Token{AccessToken: "test-access-token"}

		tokens := map[string]*oauth2.Token{
			"account_1": token,
		}

		mockCalendarService.On("Watch", anything, anything, anything, anything).Return(nil, someErr).Once()
		mockCalendarService.On("Watch", anything, anything, anything, anything).Return(mockChannel, nil)
		mockDao.On("GetUserTokens", anything, anything, anything, anything).Return(tokens["account_1"], nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockTokenSource.On("Token").Return(token, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, nil)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
			webhookURL:      "example.com",
		}

		res, errs := g.RenewWatches(ctx, watches, tokens)

		assert.NotEmpty(t, errs, "Expected  errors")
		assert.Len(t, res.Failed, 1, "Expected one watch to fail")
		assert.Equal(t, watches[0].ID, res.Failed[0].ID, "Watch IDs should match")
		assert.Len(t, res.Succesful, 1, "Expected one watch to be successfully renewed")
		assert.Equal(t, watches[0].ID, res.Succesful[0].ID, "Watch IDs should match")
		assert.Equal(t, time.UnixMilli(expirationTime), res.Succesful[0].Expiry, "Expiry should match expected time")

	})

	t.Run("error creating service", func(t *testing.T) {
		ctx := context.Background()

		mockCalendarService := mocks.NewCalendarService(t)
		mockDao := mocks2.NewDAO(t)
		mockTokenSource := mocks.NewTokenSource(t)
		mockconfig := mocks.NewOAuthConfig(t)

		watches := []*dto.WatchData{
			watch,
		}

		token := &oauth2.Token{AccessToken: "test-access-token"}

		tokens := map[string]*oauth2.Token{
			"account_1": token,
		}

		mockDao.On("GetUserTokens", anything, anything, anything, anything).Return(tokens["account_1"], nil)
		mockconfig.On("TokenSource", anything, anything).Return(mockTokenSource, nil)
		mockTokenSource.On("Token").Return(token, nil)
		mockCalendarService.On("NewService", anything, anything).Return(nil, someErr)

		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
			webhookURL:      "example.com",
		}

		res, errs := g.RenewWatches(ctx, watches, tokens)

		assert.NotEmpty(t, errs, "Expected  errors")
		assert.Len(t, res.Failed, 1, "Expected one watch to fail")
		assert.Equal(t, watches[0].ID, res.Failed[0].ID, "Watch IDs should match")

	})

	t.Run("error creating authClient", func(t *testing.T) {
		ctx := context.Background()

		mockCalendarService := mocks.NewCalendarService(t)
		mockDao := mocks2.NewDAO(t)
		mockconfig := mocks.NewOAuthConfig(t)

		watches := []*dto.WatchData{
			watch,
		}

		token := &oauth2.Token{AccessToken: "test-access-token"}

		tokens := map[string]*oauth2.Token{
			"account_1": token,
		}

		mockDao.On("GetUserTokens", anything, anything, anything, anything).Return(nil, someErr)
		g := &googleCalendar{
			dao:             mockDao,
			calendarService: mockCalendarService,
			config:          mockconfig,
			clientCache:     NewClientCache(0),
			webhookURL:      "example.com",
		}

		res, errs := g.RenewWatches(ctx, watches, tokens)

		assert.NotEmpty(t, errs, "Expected  errors")
		assert.Len(t, res.Failed, 1, "Expected one watch to fail")
		assert.Equal(t, watches[0].ID, res.Failed[0].ID, "Watch IDs should match")

	})
}
