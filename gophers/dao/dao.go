package dao

import (
	"context"
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
	SaveUserTokens(context.Context, int, string, string, time.Time) error
	GetUserTokens(context.Context, int) (*oauth2.Token, error)
	SaveUserCalendarData(context.Context, int, []*CalendarData) error
	GetUserCalendars(context.Context, int) ([]Calendar, error)
}

// New ...
func New() DAO {
	return &dao{
		DB: DB,
	}
}

type dao struct {
	DB *gorm.DB
}
