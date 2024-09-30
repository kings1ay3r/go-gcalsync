package dao

import (
	"context"
	"gcalsync/gophers/dto"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"time"
)

// DAO ...
//
//go:generate mockery --name=DAO --dir=./ --output=mocks --outpkg=mocks
type DAO interface {
	FindCalendarByCalendarID(context.Context, string) (*Calendar, error)
	InsertCalendar(context.Context, Calendar) (*Calendar, error)
	SaveUserTokens(context.Context, int, string, string, string, time.Time) error
	GetUserTokens(context.Context, int, string) (*oauth2.Token, error)
	SaveUserCalendarData(context.Context, int, []*CalendarData) error
	GetUserCalendars(context.Context, int) ([]Calendar, error)
	SaveWatch(context.Context, *Watch) error
	FindExpiringWatches(context.Context) ([]WatchesWithDetails, error)
	WatchExists(context.Context, uint, string) (bool, error)
	FindCalendarByResourceIDWithToken(context.Context, string) (*dto.CalendarDetailsByResourceIDResponse, error)
}

// New ...
func New() (DAO, error) {
	return &dao{
		DB: DB,
	}, nil
}

type dao struct {
	DB *gorm.DB
}
