package model

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string // bcrypt hash; empty for OAuth-only identities
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
