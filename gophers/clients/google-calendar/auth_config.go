package googlecalendar

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type OAuthConfigImpl struct {
	*oauth2.Config
}

func (o *OAuthConfigImpl) ClientID() string {
	return o.Config.ClientID
}

func (o *OAuthConfigImpl) GetGoogleAccountID(ctx context.Context, idToken string) (string, error) {

	payload, err := idtoken.Validate(ctx, idToken, o.Config.ClientID)
	if err != nil {
		return "", fmt.Errorf("failed to validate ID token: %w", err)
	}

	googleAccountID, ok := payload.Claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("Google Account ID (sub) not found in ID token")
	}

	return googleAccountID, nil
}

func (o *OAuthConfigImpl) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return o.Config.AuthCodeURL(state, opts...)
}

func (o *OAuthConfigImpl) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.Config.Exchange(ctx, code)
}

func (o *OAuthConfigImpl) TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource {
	return o.Config.TokenSource(ctx, t)
}
