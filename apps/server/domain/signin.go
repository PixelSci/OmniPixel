package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user is not active")
)

type User struct {
	ID                string     `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	DeviceFingerprint string     `json:"device_fingerprint,omitempty"`
	DisplayName       *string    `json:"display_name,omitempty"`
	AvatarURL         *string    `json:"avatar_url,omitempty"`
	Status            string     `json:"status"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	User        User   `json:"user"`
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	UpdateLastLogin(userID string) error
}
