package googlecalendar

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// FetchCalendars fetches the calendars for the user after handling tokens.
func (g *googleCalendar) FetchCalendars(ctx context.Context, userID int, code string) ([]*calendar.CalendarListEntry, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange auth code: %w", err)
	}

	// Save the token to the database.
	if err := g.dao.SaveUserTokens(ctx, userID, token.AccessToken, token.RefreshToken, token.Expiry); err != nil {
		return nil, fmt.Errorf("failed to save user tokens: %w", err)
	}

	tokenSource := g.config.TokenSource(ctx, token)
	client := oauth2.NewClient(ctx, tokenSource)

	// Refresh the token and save if changed
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	if newToken.AccessToken != token.AccessToken || newToken.RefreshToken != token.RefreshToken || newToken.Expiry != token.Expiry {
		if err := g.dao.SaveUserTokens(ctx, userID, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry); err != nil {
			return nil, err
		}
	}

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	calendarList, err := srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}

	return calendarList.Items, nil
}
