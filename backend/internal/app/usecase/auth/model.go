package auth

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	UserID       string         `db:"user_id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Email        string         `db:"email" json:"email"`
	Roles        pq.StringArray `db:"roles" json:"-"`
	PasswordHash []byte         `db:"password_hash" json:"-"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
}

type AuthenticatedUser struct {
	UserID      string
	Roles       []string
	Email       string
	Name        string
	DateCreated time.Time
}

type LoginUser struct {
	Email    string
	Password string
}

type DeleteToken struct {
	TokenID   string    `db:"token_id"`
	Subject   string    `db:"subject"`
	IssuedAt  time.Time `db:"issued_at"`
	ExpiresAt time.Time `db:"expires_at"`
	RevokedAt time.Time `db:"revoked_at"`
}

type ChangePassword struct {
	UserID   string
	Password string
}
