package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"omni-pixel/apperr"
	"omni-pixel/model"
)

type TokenPair struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int // seconds
}

// Claims is the JWT payload shape. Protected-route middleware (added later)
// will verify a token and hydrate a *Claims for handlers to consume.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

func NewTokenIssuer(secret string, ttl time.Duration, issuer string) *TokenIssuer {
	return &TokenIssuer{
		secret: []byte(secret),
		ttl:    ttl,
		issuer: issuer,
	}
}

func (t *TokenIssuer) Issue(u *model.User) (*TokenPair, error) {
	now := time.Now()
	claims := &Claims{
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID,
			Issuer:    t.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(t.secret)
	if err != nil {
		return nil, fmt.Errorf("sign jwt: %w", err)
	}
	return &TokenPair{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int(t.ttl.Seconds()),
	}, nil
}

// Verify parses and validates a signed token, returning the embedded claims.
// Used by the (future) auth middleware on protected routes.
func (t *TokenIssuer) Verify(raw string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(raw, &Claims{}, func(tok *jwt.Token) (any, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}
		return t.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperr.ErrTokenExpired.WithCause(err)
		}
		return nil, apperr.ErrTokenInvalid.WithCause(err)
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, apperr.ErrTokenInvalid
	}
	return claims, nil
}
