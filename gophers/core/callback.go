package core

import (
	"context"
	"fmt"
	"gcalsync/gophers/clients/logger"
	"gcalsync/gophers/dao"
	"gcalsync/gophers/middlewares/auth"
	"google.golang.org/api/calendar/v3"
	"time"
)

// InsertCalendars ...
func (c *calendarClient) InsertCalendars(ctx context.Context, code string) error {

	log := logger.GetInstance()
	currUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	currUserID := currUser.ID
	calendars, accountID, err := c.googleCalClient.FetchCalendars(ctx, currUserID, code)

	if err != nil {
		return fmt.Errorf("unable to retrieve calendar list: %w", err)
	}

	var dbCalendars []*dao.CalendarData
	eventLength := 0
	for _, calendarEntry := range calendars {
		dbCalendar := &dao.CalendarData{
			CalendarID: calendarEntry.Id,
			Name:       calendarEntry.Summary,
			AccountID:  accountID,
		}
		events, err := c.googleCalClient.FetchEventsWithUserID(ctx, currUserID, accountID, calendarEntry.Id)

		if err != nil {
			// TODO: Implement a dead letter queue for error handling
			log.Error(ctx, "unable to retrieve calendar list: %v", err)
			continue
		}
		for _, event := range events {
			dbEvent, err := mapEvent(event)
			eventLength++
			if err != nil {
				log.Error(ctx, "failed to map event: %v", err)
			}
			dbCalendar.Events = append(dbCalendar.Events, dbEvent)
		}
		dbCalendars = append(dbCalendars, dbCalendar)
	}
	log.Info(ctx, "Logging %v events across %v calendars", eventLength, len(dbCalendars))

	go c.insertInBackground(ctx, currUserID, dbCalendars)

	return nil
}

func (c *calendarClient) insertInBackground(ctx context.Context, userID int, dbCalendars []*dao.CalendarData) {
	// TODO: Ensure Semaphore lock to avoid issues with concurrency
	err := c.dao.SaveUserCalendarData(ctx, userID, dbCalendars)
	if err != nil {
		logger.GetInstance().Error(ctx, "unable to insert calendar list: %v", err)
	}
	return
}

func mapEvent(e *calendar.Event) (dao.EventData, error) {
	var startTime, endTime time.Time
	var err error

	// TODO: Extract helper method and add unit tests to handle time. Handle Timezone
	// Handle start time
	if e.Start.DateTime != "" {
		startTime, err = time.Parse(time.RFC3339Nano, e.Start.DateTime)
		if err != nil {
			return dao.EventData{}, err
		}
	} else if e.Start.Date != "" {
		startTime, err = time.Parse("2006-01-02", e.Start.Date)
		if err != nil {
			return dao.EventData{}, err
		}
	}
	// Handle end time
	if e.End.DateTime != "" {
		endTime, err = time.Parse(time.RFC3339Nano, e.End.DateTime)
		if err != nil {
			return dao.EventData{}, err
		}
	} else if e.End.Date != "" {
		endTime, err = time.Parse("2006-01-02", e.End.Date)
		if err != nil {
			return dao.EventData{}, err
		}
	}

	// If end date is provided without a time, assume end time is the end of the day
	if e.End.Date != "" && e.End.DateTime == "" {
		// FIXME: Convert to GMT, compute end time
		endTime = endTime.Add(24 * time.Hour).Add(-time.Nanosecond) // Set end time to 23:59:59.999999999
	}

	if startTime.IsZero() {
		return dao.EventData{}, fmt.Errorf("event %s has no valid start time", e.Id)
	}
	if endTime.IsZero() {
		return dao.EventData{}, fmt.Errorf("event %s has no valid end time", e.Id)
	}

	return dao.EventData{
		EventID:   e.Id,
		Name:      e.Summary,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}
