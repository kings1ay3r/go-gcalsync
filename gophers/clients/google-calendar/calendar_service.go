package googlecalendar

import (
	"context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"net/http"
)

type googleCalendarService struct {
}

// TODO: Implement Retries / Fault accommodation for requests

func (g *googleCalendarService) NewService(ctx context.Context, client *http.Client) (*calendar.Service, error) {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return srv, err
	}
	return srv, nil
}

func (g *googleCalendarService) ListEvents(ctx context.Context, srv *calendar.Service, calendarID string) (*calendar.Events, error) {
	return srv.Events.List(calendarID).Do()
}

func (g *googleCalendarService) ListCalendars(ctx context.Context, srv *calendar.Service) (*calendar.CalendarList, error) {
	return srv.CalendarList.List().Do()
}
