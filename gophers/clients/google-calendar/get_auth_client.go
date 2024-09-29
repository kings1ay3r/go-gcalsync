package googlecalendar

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func (g *googleCalendar) getOAuthClientForUser(ctx context.Context, token *oauth2.Token, userID int, accountID string) (*http.Client, error) {

	// There is an assumption here that tokens expire in one hour
	//  and cache expired in less than one hour. This prevents stale clients in cache.
	if client, ok := g.clientCache.Get(accountID); ok {
		return client, nil
	}
	if token == nil || token.AccessToken == "" || token.RefreshToken == "" {
		var err error
		token, err = g.dao.GetUserTokens(ctx, userID, accountID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user tokens: %w", err)
		}
	}

	if token.Expiry.IsZero() || token.Expiry.Sub(time.Now()) < 30*time.Second {
		tokenSource := g.config.TokenSource(ctx, token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh tokens: %w", err)
		}

		if newToken.AccessToken != token.AccessToken || newToken.RefreshToken != token.RefreshToken || newToken.Expiry != token.Expiry {
			if err := g.dao.SaveUserTokens(ctx, userID, accountID, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry); err != nil {
				return nil, fmt.Errorf("failed to save refreshed token: %w", err)
			}
		}
		token = newToken
	}

	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	g.clientCache.Push(accountID, client)
	return client, nil
}
