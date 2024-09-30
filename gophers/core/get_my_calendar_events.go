package core

import (
	"context"
	"gcalsync/gophers/dto"
	"gcalsync/gophers/middlewares/auth"
)

// GetMyCalendarEvents ...
func (c *calendarClient) GetMyCalendarEvents(ctx context.Context) ([]dto.Calendar, error) {
	currUser, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	calendars, err := c.dao.GetUserCalendars(ctx, currUser.ID)
	if err != nil {
		return nil, err
	}
	var resp []dto.Calendar
	for _, calendar := range calendars {

		var dtoEvents []dto.Event
		for _, event := range calendar.Events {
			dtoEvent := dto.Event{
				ID:      event.EventID,
				Summary: event.Summary,
				Start:   event.StartTime,
				End:     event.EndTime,
			}
			dtoEvents = append(dtoEvents, dtoEvent)
		}

		resp = append(resp, dto.Calendar{
			ID:     calendar.CalendarID,
			Name:   calendar.Name,
			Events: dtoEvents,
		})

	}
	return resp, nil
}
