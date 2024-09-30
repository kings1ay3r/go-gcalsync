package core

import (
	"context"
	googlecalendar "gcalsync/gophers/clients/google-calendar"
	"gcalsync/gophers/clients/logger"
	"gcalsync/gophers/dao"
	"gcalsync/gophers/dto"
)

//go:generate mockery --name=Core --dir=./ --output=mocks --outpkg=mocks
type Core interface {
	InsertCalendars(ctx context.Context, code string) error
	GetAuthCodeURL(ctx context.Context) (string, error)
	GetMyCalendarEvents(ctx context.Context) ([]dto.Calendar, error)
	RenewExpiringWatches(ctx context.Context)
}

func New() (Core, error) {

	client, err := googlecalendar.New()

	if err != nil {
		logger.GetInstance().Error(nil, "unable to init services : %v", err)
		return nil, err
	}
	daoInst, err := dao.New()
	if err != nil {
		return nil, err
	}
	return &calendarClient{
		googleCalClient: client,
		dao:             daoInst,
	}, nil
}

type calendarClient struct {
	googleCalClient googlecalendar.GoogleCalendar
	dao             dao.DAO
}
