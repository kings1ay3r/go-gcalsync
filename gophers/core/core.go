package core

import (
	"context"
	googlecalendar "gcalsync/gophers/clients/google-calendar"
	"gcalsync/gophers/clients/logger"
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

func New() (Core, error) {

	client, err := googlecalendar.New()

	if err != nil {
		logger.GetInstance().Error(nil, "unable to init services : %v", err)
		return nil, err
	}
	dao, err := dao.New()
	if err != nil {
		return nil, err
	}
	return &calendarClient{
		googleCalClient: client,
		dao:             dao,
	}, nil
}

type calendarClient struct {
	googleCalClient googlecalendar.GoogleCalendar
	dao             dao.DAO
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
