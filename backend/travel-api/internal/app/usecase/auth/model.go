package auth

import (
	"time"

	"github.com/lib/pq"
)

// TODO: move to db layer
type User struct {
	UserID       string         `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	IsActive     bool           `db:"is_active"`
	TokenVersion int32          `db:"token_version"`
	Roles        pq.StringArray `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

type AuthenticatedUser struct {
	UserID       string
	Roles        []string
	Email        string
	Name         string
	TokenVersion int32
	DateCreated  time.Time
}

type LoginUser struct {
	Email    string
	Password string
}

type DeleteToken struct {
	TokenID      string    `db:"token_id"`
	Subject      string    `db:"subject"`
	TokenVersion int32     `db:"token_version"`
	IssuedAt     time.Time `db:"issued_at"`
	ExpiresAt    time.Time `db:"expires_at"`
	RevokedAt    time.Time `db:"revoked_at"`
}

type ChangePassword struct {
	UserID      string
	Password    string
	PasswordOld string
}

type ResetPassword struct {
	Email      string
	Name       string
	ResetToken string
}

type ResetToken struct {
	TokenID   string    `db:"token_id"`
	UserID    string    `db:"user_id"`
	Email     string    `db:"email"`
	ExpiresAt time.Time `db:"expires_at"`
	IssuedAt  time.Time `db:"issued_at"`
}

type SubmitPassword struct {
	ResetToken string
	Password   string
}
