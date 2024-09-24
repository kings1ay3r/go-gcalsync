package google_calendar

import (
	"context"

	"gcalsync/gophers/dto"

	"google.golang.org/api/calendar/v3"
)

// GoogleCalendarDAO implements the CalendarService interface.
type GoogleCalendarDAO struct {
	service *calendar.Service
}

// NewGoogleCalendarDAO creates a new instance of GoogleCalendarDAO.
func NewGoogleCalendarDAO(service *calendar.Service) dto.CalendarService {
	return &GoogleCalendarDAO{service: service}
}

// ListEvents fetches events from the Google Calendar API.
func (g *GoogleCalendarDAO) ListEvents(ctx context.Context, calendarID string) ([]dto.Event, error) {
	events, err := g.service.Events.List(calendarID).Do()
	if err != nil {
		return nil, err
	}

	var eventList []dto.Event
	for _, e := range events.Items {
		eventList = append(eventList, dto.Event{
			ID:      e.Id,
			Summary: e.Summary,
			Start:   e.Start.DateTime,
			End:     e.End.DateTime,
		})
	}
	return eventList, nil
}
