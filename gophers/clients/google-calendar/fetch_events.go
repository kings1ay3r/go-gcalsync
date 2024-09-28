package googlecalendar

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

// fetchEvents fetches events for a given calendar using the provided token.
func (g *googleCalendar) fetchEvents(ctx context.Context, calendarID string, token *oauth2.Token) ([]*calendar.Event, error) {
	tokenSource := g.config.TokenSource(ctx, token)
	client := oauth2.NewClient(ctx, tokenSource)

	// Refresh the token and save if necessary.
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	if newToken.AccessToken != token.AccessToken || newToken.RefreshToken != token.RefreshToken || newToken.Expiry != token.Expiry {
		userID, ok := ctx.Value("currentUser").(int)
		if !ok {
			return nil, fmt.Errorf("could not get current user id")
		}
		if err := g.dao.SaveUserTokens(ctx, userID, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry); err != nil {
			return nil, err
		}
	}

	err = g.calendarService.NewService(ctx, client)
	if err != nil {
		return nil, err
	}

	eventsList, err := g.calendarService.ListEvents(ctx, calendarID)
	if err != nil {
		return nil, err
	}

	return eventsList.Items, nil
}

// FetchEventsWithCode exchanges the authorization code for a token, saves the token to the database, and fetches events.
func (g *googleCalendar) FetchEventsWithCode(ctx context.Context, userID int, code string, calendarID string) ([]*calendar.Event, error) {
	// Exchange the authorization code for an access token.
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange auth code: %w", err)
	}

	// Save the token to the database.
	if err := g.dao.SaveUserTokens(ctx, userID, token.AccessToken, token.RefreshToken, token.Expiry); err != nil {
		return nil, fmt.Errorf("failed to save user tokens: %w", err)
	}

	// Fetch events using the saved token.
	return g.fetchEvents(ctx, calendarID, token)
}

// FetchEventsWithUserID retrieves tokens from the database and fetches events.
func (g *googleCalendar) FetchEventsWithUserID(ctx context.Context, userID int, calendarID string) ([]*calendar.Event, error) {
	// Retrieve the tokens from the database.
	token, err := g.dao.GetUserTokens(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	// Fetch events using the retrieved token.
	return g.fetchEvents(ctx, calendarID, token)
}
