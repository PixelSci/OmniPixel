// Package internal holds built-in helpers shared by the app (tokens, crypto wrappers, etc.).
// It is not where Router / Controller / Usecase / Repository live — those sit at the module root.
package internal

import "omni-pixel/domain"

// CreateAccessToken issues a JWT access token for the given user (stub — implement with your JWT library).
func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return "", nil
}
