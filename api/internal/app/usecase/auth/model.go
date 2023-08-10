package auth

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	// TODO: validation
	ID           string         `db:"user_id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Email        string         `db:"email" json:"email"`
	Roles        pq.StringArray `db:"roles" json:"roles"`
	PasswordHash []byte         `db:"password_hash" json:"-"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
}

type AuthenticatedUser struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
	Email string   `json:"email"`
	Name  string   `json:"name"`
}

// TODO: validation
type LoginUser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Login struct {
	// TODO: validation
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TODO: validation
type Logout struct {
	Email string `json:"email" validate:"required"`
}

// TODO: validation
type Refresh struct {
	Token string `json:"token" validate:"required"`
}

// TODO: validation
type PasswordReset struct {
	Email string `json:"email" validate:"required"`
}

type StoreToken struct {
	Token     string `db:"token"`
	UserID    string `db:"user_id"`
	ExpiresAt int64  `db:"expires_at"`
	IssuedAt  int64  `db:"issued_at"`
	CreatedAt int64  `db:"created_at"`
}
