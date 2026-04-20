// Package auth owns identity — registration, credential verification, and
// access-token issuing. HTTP handlers depend on *Service; provider plugins
// (password today, github/google later) live beside it and implement one of
// the narrow interfaces below.
package auth

import (
	"context"

	"omni-pixel/model"
)

// Provider is the common identity of any auth plugin.
type Provider interface {
	Name() string
}

// PasswordAuthenticator verifies an email + password pair against stored
// credentials. Exactly one built-in implementation today (see password.go).
type PasswordAuthenticator interface {
	Provider
	Authenticate(ctx context.Context, email, password string) (*model.User, error)
}

// OAuthAuthenticator exchanges a third-party OAuth2 authorization code for a
// local user. GitHub and Google will implement this when they're added —
// each as its own file in this package (e.g. github.go, google.go) and
// registered on *Service alongside the password authenticator.
type OAuthAuthenticator interface {
	Provider
	AuthCodeURL(state string) string
	Exchange(ctx context.Context, code string) (*model.User, error)
}
