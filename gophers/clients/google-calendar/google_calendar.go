package googlecalendar

import (
	"context"
	"errors"
	"gcalsync/gophers/dao"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"net/url"
	"os"
)

// GoogleCalendar ...
//
//go:generate mockery --name=GoogleCalendar --dir=./ --output=mocks --outpkg=mocks
type GoogleCalendar interface {
	FetchCalendars(ctx context.Context, userID int, code string) ([]*calendar.CalendarListEntry, error)
	FetchEventsWithCode(ctx context.Context, userID int, code string, calendarID string) ([]*calendar.Event, error)
	FetchEventsWithUserID(ctx context.Context, userID int, calendarID string) ([]*calendar.Event, error)
	GetAuthCodeURL(ctx context.Context, userToken string) string
}

type googleCalendar struct {
	config *oauth2.Config
	dao    dao.DAO
}

// New initializes a GoogleCalendar instance with the DAO.
func New() (GoogleCalendar, error) {
	daoInstance := dao.New()
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET not set")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarScope},
		RedirectURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
	}

	return &googleCalendar{
		config: config,
		dao:    daoInstance, // Pass DAO instance here
	}, nil
}

// GetAuthCodeURL ...
func (g *googleCalendar) GetAuthCodeURL(_ context.Context, userToken string) string {
	encodedUserID := url.QueryEscape(userToken)
	return g.config.AuthCodeURL(encodedUserID, oauth2.AccessTypeOffline)
}
