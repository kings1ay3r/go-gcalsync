package core

import (
	"context"
	"fmt"
	"gcalsync/gophers/middlewares/auth"
)

// GetAuthCodeURL returns url to redirect to for authorization
func (c *calendarClient) GetAuthCodeURL(ctx context.Context) (string, error) {

	userKey, ok := ctx.Value(auth.ContextUserIDKey).(string)
	if !ok {
		return "", fmt.Errorf("failed to get auth code url")
	}
	return c.googleCalClient.GetAuthCodeURL(ctx, userKey), nil
}
