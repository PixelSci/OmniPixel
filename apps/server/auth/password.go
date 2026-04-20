package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"omni-pixel/apperr"
	"omni-pixel/model"
	"omni-pixel/repository"
)

// hashPassword produces a bcrypt hash suitable for storing in users.password_hash.
func hashPassword(plain string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// verifyPassword returns nil when plain matches the stored bcrypt hash.
func verifyPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

// NewPasswordProvider returns the default built-in credential checker.
func NewPasswordProvider(users repository.UserRepository) PasswordAuthenticator {
	return &passwordProvider{users: users}
}

type passwordProvider struct {
	users repository.UserRepository
}

func (p *passwordProvider) Name() string { return "password" }

func (p *passwordProvider) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	user, err := p.users.FindByEmail(ctx, email)
	if err != nil {
		// Return the same error for missing user and wrong password so
		// callers can't probe for registered emails.
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperr.ErrPasswordWrong
		}
		return nil, apperr.ErrInternal.WithCause(err)
	}
	if user.PasswordHash == "" {
		// User registered via OAuth only; no password set.
		return nil, apperr.ErrPasswordWrong
	}
	if err := verifyPassword(user.PasswordHash, password); err != nil {
		return nil, apperr.ErrPasswordWrong
	}
	return user, nil
}
