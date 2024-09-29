package googlecalendar

import (
	"context"
	"fmt"
	"google.golang.org/api/calendar/v3"
)

// FetchCalendars fetches the calendars for the user after handling tokens.
func (g *googleCalendar) FetchCalendars(ctx context.Context, userID int, code string) ([]*calendar.CalendarListEntry, string, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, "", err
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, "", fmt.Errorf("no id_token field in oauth2 token")
	}
	googleID, err := g.config.GetGoogleAccountID(ctx, idToken)
	if err != nil {
		return nil, "", fmt.Errorf("google id not found: %w", err)
	}

	if err := g.dao.SaveUserTokens(ctx, userID, googleID, token.AccessToken, token.RefreshToken, token.Expiry); err != nil {
		return nil, "", fmt.Errorf("failed to save user tokens: %w", err)
	}

	client, err := g.getOAuthClientForUser(ctx, token, userID, googleID)
	if err != nil {
		return nil, "", err
	}
	srv, err := g.calendarService.NewService(ctx, client)
	if err != nil {
		return nil, "", err
	}
	calendarList, err := g.calendarService.ListCalendars(ctx, srv)
	if err != nil {
		return nil, "", err
	}

	return calendarList.Items, googleID, nil
}
