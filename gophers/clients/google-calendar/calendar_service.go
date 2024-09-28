package googlecalendar

import (
	"context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"net/http"
)

type googleCalendarService struct {
	srv *calendar.Service
}

func (g *googleCalendarService) NewService(ctx context.Context, client *http.Client) error {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	g.srv = srv
	return nil
}

func (g *googleCalendarService) ListEvents(ctx context.Context, calendarID string) (*calendar.Events, error) {
	return g.srv.Events.List(calendarID).Do()
}

func (g *googleCalendarService) ListCalendars(ctx context.Context) (*calendar.CalendarList, error) {
	return g.srv.CalendarList.List().Do()
}
