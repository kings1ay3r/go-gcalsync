package core

import (
	"context"
	"fmt"
	googlecalendar "gcalsync/gophers/clients/google-calendar"
	"gcalsync/gophers/dao"
	"gcalsync/gophers/dto"
	"gcalsync/gophers/middlewares/auth"
)

//go:generate mockery --name=Core --dir=./ --output=mocks --outpkg=mocks
type Core interface {
	InsertCalendars(ctx context.Context, code string) error
	GetAuthCodeURL(ctx context.Context) (string, error)
	GetMyCalendarEvents(ctx context.Context) ([]dto.Calendar, error)
}

func New() Core {

	client, err := googlecalendar.New()
	if err != nil {
		panic(err)
	}
	return &calendarClient{
		googleCalClient: client,
		dao:             dao.New(),
	}
}

type calendarClient struct {
	googleCalClient googlecalendar.GoogleCalendar
	dao             dao.DAO
}

// GetAuthCodeURL returns url to redirect to for authorization
func (c *calendarClient) GetAuthCodeURL(ctx context.Context) (string, error) {

	userKey, ok := ctx.Value(auth.ContextUserIDKey).(string)
	if !ok {
		return "", fmt.Errorf("failed to get auth code url")
	}
	return c.googleCalClient.GetAuthCodeURL(ctx, userKey), nil
}

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
