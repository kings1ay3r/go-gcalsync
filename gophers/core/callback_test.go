package core

import (
	"context"
	"fmt"
	gcmocks "gcalsync/gophers/clients/google-calendar/mocks"
	"gcalsync/gophers/middlewares/auth"
	"github.com/stretchr/testify/mock"
	"time"

	daomocks "gcalsync/gophers/dao/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
	"testing"
)

func TestInsertCalendars(t *testing.T) {
	anything := mock.Anything

	t.Run("Successfully inserts calendars and events", func(t *testing.T) {
		ctx := getContextWithSession()
		mockDAO := daomocks.NewDAO(t)
		mockGoogleCalClient := gcmocks.NewGoogleCalendar(t)
		calendarClient := calendarClient{
			googleCalClient: mockGoogleCalClient,
			dao:             mockDAO,
		}
		mockCalendars := []*calendar.CalendarListEntry{
			{Id: "cal1", Summary: "Test Calendar 1"},
		}

		mockEvents := []*calendar.Event{
			{
				Id:      "event1",
				Summary: "Test Event 1",
				Start:   &calendar.EventDateTime{DateTime: "2023-09-28T10:00:00Z"},
				End:     &calendar.EventDateTime{DateTime: "2023-09-28T11:00:00Z"},
			},
		}

		mockGoogleCalClient.On("FetchCalendars", anything, anything, anything).Return(mockCalendars, "", nil)
		mockDAO.On("SaveUserCalendarData", anything, anything, anything).Return(nil)
		mockGoogleCalClient.On("FetchEventsWithUserID", anything, anything, anything, anything).Return(mockEvents, nil)

		err := calendarClient.InsertCalendars(ctx, "mock-code")

		// Allow the goroutine to run
		time.Sleep(100 * time.Millisecond)

		assert.NoError(t, err)
		mockGoogleCalClient.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
	})

	t.Run("Error when fetching user ID from context", func(t *testing.T) {
		ctx := context.Background()
		calendarClient := calendarClient{}

		err := calendarClient.InsertCalendars(ctx, "mock-code")

		assert.Error(t, err)
		assert.Equal(t, "missing user context", err.Error())
	})

	t.Run("Error when fetching calendars from Google Calendar API", func(t *testing.T) {
		ctx := getContextWithSession()
		mockDAO := daomocks.NewDAO(t)
		mockGoogleCalClient := gcmocks.NewGoogleCalendar(t)
		calendarClient := calendarClient{
			googleCalClient: mockGoogleCalClient,
			dao:             mockDAO,
		}
		mockGoogleCalClient.On("FetchCalendars", anything, anything, anything).Return(nil, "", fmt.Errorf("unable to fetch calendars"))

		err := calendarClient.InsertCalendars(ctx, "mock-code")

		assert.Error(t, err)
		assert.Equal(t, "unable to retrieve calendar list: unable to fetch calendars", err.Error())
		mockDAO.AssertNotCalled(t, "SaveUserCalendarData", ctx, anything, anything)
	})

	t.Run("Execution continues when fetching events from Google Calendar API fails", func(t *testing.T) {
		ctx := getContextWithSession()
		mockDAO := daomocks.NewDAO(t)
		mockGoogleCalClient := gcmocks.NewGoogleCalendar(t)
		calendarClient := calendarClient{
			googleCalClient: mockGoogleCalClient,
			dao:             mockDAO,
		}
		mockCalendars := []*calendar.CalendarListEntry{
			{Id: "cal1", Summary: "Test Calendar 1"},
		}
		mockGoogleCalClient.On("FetchCalendars", anything, anything, anything).Return(mockCalendars, "", nil)
		mockGoogleCalClient.On("FetchEventsWithUserID", anything, anything, anything, anything).Return(nil, fmt.Errorf("unable to fetch events"))
		mockDAO.On("SaveUserCalendarData", anything, anything, anything).Return(nil)

		err := calendarClient.InsertCalendars(ctx, "mock-code")

		// Allow the goroutine to run
		time.Sleep(100 * time.Millisecond)

		assert.NoError(t, err)
	})

	t.Run("Error when saving calendar data in background", func(t *testing.T) {

		ctx := getContextWithSession()
		mockDAO := daomocks.NewDAO(t)
		mockGoogleCalClient := gcmocks.NewGoogleCalendar(t)
		calendarClient := calendarClient{
			googleCalClient: mockGoogleCalClient,
			dao:             mockDAO,
		}
		mockCalendars := []*calendar.CalendarListEntry{
			{Id: "cal1", Summary: "Test Calendar 1"},
		}
		mockEvents := []*calendar.Event{
			{
				Id:      "event1",
				Summary: "Test Event 1",
				Start:   &calendar.EventDateTime{DateTime: "2023-09-28T10:00:00Z"},
				End:     &calendar.EventDateTime{DateTime: "2023-09-28T11:00:00Z"},
			},
		}
		mockGoogleCalClient.On("FetchCalendars", anything, anything, anything).Return(mockCalendars, "", nil)
		mockGoogleCalClient.On("FetchEventsWithUserID", anything, anything, anything, anything).Return(mockEvents, nil)
		mockDAO.On("SaveUserCalendarData", anything, anything, anything).Return(fmt.Errorf("unable to save calendar data"))

		err := calendarClient.InsertCalendars(ctx, "mock-code")

		// Allow the goroutine to run
		time.Sleep(100 * time.Millisecond)

		assert.NoError(t, err)
		mockDAO.AssertExpectations(t)

	})

	t.Run("Error when saving calendars to the database", func(t *testing.T) {
		ctx := getContextWithSession()
		mockDAO := daomocks.NewDAO(t)
		mockGoogleCalClient := gcmocks.NewGoogleCalendar(t)
		calendarClient := calendarClient{
			googleCalClient: mockGoogleCalClient,
			dao:             mockDAO,
		}
		mockCalendars := []*calendar.CalendarListEntry{
			{Id: "cal1", Summary: "Test Calendar 1"},
		}
		mockEvents := []*calendar.Event{
			{
				Id:      "event1",
				Summary: "Test Event 1",
				Start:   &calendar.EventDateTime{DateTime: "2023-09-28T10:00:00Z"},
				End:     &calendar.EventDateTime{DateTime: "2023-09-28T11:00:00Z"},
			},
		}
		mockGoogleCalClient.On("FetchCalendars", anything, anything, anything).Return(mockCalendars, "", nil)
		mockGoogleCalClient.On("FetchEventsWithUserID", anything, anything, anything, anything).Return(mockEvents, nil)
		mockDAO.On("SaveUserCalendarData", anything, anything, anything).Return(fmt.Errorf("database error"))

		err := calendarClient.InsertCalendars(ctx, "mock-code")

		// Allow the goroutine to run
		time.Sleep(100 * time.Millisecond)

		assert.NoError(t, err)
		mockGoogleCalClient.AssertExpectations(t)
		mockDAO.AssertExpectations(t)
	})
}
func getContextWithSession() context.Context {
	mockUser := auth.Session{ID: 123}
	ctx := context.Background()
	ctx = context.WithValue(ctx, auth.ContextUserKey, mockUser)
	ctx = context.WithValue(ctx, auth.ContextUserIDKey, "ID123")
	return ctx
}
