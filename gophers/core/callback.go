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

	currUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	currUserID := currUser.ID
	calendars, err := c.googleCalClient.FetchCalendars(ctx, currUserID, code)
	if err != nil {
		return fmt.Errorf("unable to retrieve calendar list: %w", err)
	}

	var dbCalendars []*dao.CalendarData
	eventLength := 0
	for _, calendarEntry := range calendars {
		dbCalendar := &dao.CalendarData{
			CalendarID: calendarEntry.Id,
			Name:       calendarEntry.Summary,
		}
		events, err := c.googleCalClient.FetchEventsWithUserID(ctx, currUserID, calendarEntry.Id)

		if err != nil {
			return fmt.Errorf("unable to retrieve calendar list: %w", err)
		}
		for _, event := range events {
			dbEvent, err := mapEvent(event)
			eventLength++
			if err != nil {
				fmt.Println(err)
			}
			dbCalendar.Events = append(dbCalendar.Events, dbEvent)
		}
		dbCalendars = append(dbCalendars, dbCalendar)
	}
	logger.GetInstance().Info(ctx, "Logging %v events across %v calendars", eventLength, len(dbCalendars))
	go c.insertInBackground(ctx, err, dbCalendars)

	return nil
}

func (c *calendarClient) insertInBackground(ctx context.Context, err error, dbCalendars []*dao.CalendarData) {
	// FIXME: Get and insert current user from ctx
	err = c.dao.SaveUserCalendarData(ctx, 1, dbCalendars)
	if err != nil {
		logger.GetInstance().Error(ctx, "unable to insert calendar list: %v", err)
	}
	return
}

func mapEvent(e *calendar.Event) (dao.EventData, error) {
	var startTime, endTime time.Time
	var err error

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
		endTime = endTime.Add(24 * time.Hour).Add(-time.Nanosecond) // Set end time to 23:59:59.999999999
	}

	return dao.EventData{
		EventID:   e.Id,
		Name:      e.Summary,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}
