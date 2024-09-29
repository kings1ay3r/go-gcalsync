package googlecalendar

import (
	"context"
	"errors"
	"fmt"
	"gcalsync/gophers/dao"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GoogleCalendar ...
//
//go:generate mockery --name=GoogleCalendar --dir=./ --output=mocks --outpkg=mocks
type GoogleCalendar interface {
	FetchCalendars(context.Context, int, string) ([]*calendar.CalendarListEntry, string, error)
	FetchEventsWithCode(context.Context, int, string, string, string) ([]*calendar.Event, error)
	FetchEventsWithUserID(context.Context, int, string, string) ([]*calendar.Event, error)
	GetAuthCodeURL(context.Context, string) string
}

//go:generate mockery --name=OAuthConfig --dir=./ --output=mocks --outpkg=mocks
type OAuthConfig interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
	ClientID() string
	GetGoogleAccountID(context.Context, string) (string, error)
}

//go:generate mockery --name=CalendarService --dir=./ --output=mocks --outpkg=mocks
type CalendarService interface {
	ListEvents(context.Context, *calendar.Service, string) (*calendar.Events, error)
	ListCalendars(context.Context, *calendar.Service) (*calendar.CalendarList, error)
	NewService(context.Context, *http.Client) (*calendar.Service, error)
}

//go:generate mockery --name=TokenSource --dir=./ --output=mocks --outpkg=mocks
type TokenSource interface {
	Token() (*oauth2.Token, error)
}

var config *oauth2.Config

type googleCalendar struct {
	config          OAuthConfig
	dao             dao.DAO
	calendarService CalendarService
	clientCache     *ClientCache
}

// New initializes a GoogleCalendar instance with the DAO.
func New() (GoogleCalendar, error) {
	daoInstance, err := dao.New()
	if err != nil {
		return nil, fmt.Errorf("failed to init dao: %w", err)
	}
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, errors.New("GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, or GOOGLE_CALLBACK_URL not set")
	}
	config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarScope, "email", "profile"},
		RedirectURL:  redirectURL,
	}

	return &googleCalendar{
		config:          &OAuthConfigImpl{config},
		dao:             daoInstance,
		calendarService: &googleCalendarService{},
		clientCache:     NewClientCache(time.Hour - time.Minute*5), //Keeping expiry at 55 minutes, as google expiry is per hour
	}, nil
}

// GetAuthCodeURL ...
func (g *googleCalendar) GetAuthCodeURL(_ context.Context, userToken string) string {
	// TODO: Setup JWT Token / Token User Mapping for security
	encodedUserID := url.QueryEscape(userToken)
	return g.config.AuthCodeURL(encodedUserID, oauth2.AccessTypeOffline)
}
