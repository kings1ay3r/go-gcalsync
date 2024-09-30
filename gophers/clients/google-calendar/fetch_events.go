package googlecalendar

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

// fetchEvents fetches events for a given calendar using the provided token.
func (g *googleCalendar) fetchEvents(ctx context.Context, calendarID string, token *oauth2.Token, userID int, accountID string) ([]*calendar.Event, error) {

	client, err := g.getOAuthClientForUser(ctx, token, userID, accountID)
	if err != nil {
		return nil, err
	}
	srv, err := g.calendarService.NewService(ctx, client)
	if err != nil {
		return nil, err
	}

	eventsList, err := g.calendarService.ListEvents(ctx, srv, calendarID)

	// TODO: Implement Use of NextSyncToken to get Delta
	// TODO: Implement ListOptions to get deletedEvents
	// TODO: Store Event Decorators
	if err != nil {
		return nil, err
	}

	return eventsList.Items, nil
}

// FetchEventsWithCode exchanges the authorization code for a token, saves the token to the database, and fetches events.
func (g *googleCalendar) FetchEventsWithCode(ctx context.Context, userID int, code string, accountID string, calendarID string) ([]*calendar.Event, error) {
	// Exchange the authorization code for an access token.
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange auth code: %w", err)
	}

	// Save the token to the database.
	if err := g.dao.SaveUserTokens(ctx, userID, accountID, token.AccessToken, token.RefreshToken, token.Expiry); err != nil {
		return nil, fmt.Errorf("failed to save user tokens: %w", err)
	}

	// Fetch events using the saved token.
	return g.fetchEvents(ctx, calendarID, token, userID, accountID)
}

// FetchEventsWithUserID retrieves tokens from the database and fetches events.
func (g *googleCalendar) FetchEventsWithUserID(ctx context.Context, userID int, accountID string, calendarID string) ([]*calendar.Event, error) {
	return g.fetchEvents(ctx, calendarID, nil, userID, accountID)
}

// FetchEventsFromResource ...
func (g *googleCalendar) FetchEventsFromResource(ctx context.Context, resourceID string) ([]*calendar.Event, int, error) {
	watch, err := g.dao.FindCalendarByResourceIDWithToken(ctx, resourceID)
	if err != nil {
		return nil, 0, err
	}
	events, err := g.fetchEvents(
		ctx,
		watch.GoogleCalendarID,
		nil,
		watch.UserID,
		watch.AccountID,
	)
	return events, watch.UserID, nil
}
