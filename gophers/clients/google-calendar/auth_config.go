package googlecalendar

import (
	"context"
	"golang.org/x/oauth2"
)

type OAuthConfigImpl struct {
	*oauth2.Config
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
