// Package internal holds built-in helpers shared by the app (tokens, crypto wrappers, etc.).
// It is not where Router / Controller / Usecase / Repository live — those sit at the module root.
package internal

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"omni-pixel/domain"
)

func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	now := time.Now()
	claims := domain.JwtCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiry) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
