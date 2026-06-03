package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user is not active")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type User struct {
	ID                string     `json:"id" gorm:"primaryKey;column:id"`
	Username          string     `json:"username" gorm:"column:username"`
	Email             string     `json:"email" gorm:"column:email"`
	PasswordHash      string     `json:"-" gorm:"column:password_hash"`
	DeviceFingerprint string     `json:"device_fingerprint,omitempty" gorm:"column:device_fingerprint"`
	DisplayName       *string    `json:"display_name,omitempty" gorm:"column:display_name"`
	AvatarURL         *string    `json:"avatar_url,omitempty" gorm:"column:avatar_url"`
	Status            string     `json:"status" gorm:"column:status;default:active"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty" gorm:"column:email_verified_at"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty" gorm:"column:last_login_at"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"column:updated_at"`
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

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	User        User   `json:"user"`
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	ExistsByEmail(email string) error
	Create(user *User) error
	UpdateLastLogin(userID string) error
}
