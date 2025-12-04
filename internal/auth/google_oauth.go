package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"qr-saas/internal/config"
)

type GoogleOAuth struct {
	cfg    config.Config
	oauth2 *oauth2.Config
}

func NewGoogleOAuth(cfg config.Config) *GoogleOAuth {
	return &GoogleOAuth{
		cfg: cfg,
		oauth2: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *GoogleOAuth) AuthCodeURL(state string) string {
	return g.oauth2.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

type GoogleUserInfo struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified_email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
	Sub        string `json:"sub"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

func (g *GoogleOAuth) ExchangeAndFetchUser(ctx context.Context, code string) (*GoogleUserInfo, error) {
	token, err := g.oauth2.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := g.oauth2.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo status: %d", resp.StatusCode)
	}

	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
