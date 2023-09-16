package user

import (
	"time"

	"github.com/lib/pq"
)

// TODO: refactor to db layer
type User struct {
	ID           string         `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	IsActive     bool           `db:"is_active"`
	IsDeleted    bool           `db:"is_deleted"`
	TokenVersion int32          `db:"token_version"`
	Roles        pq.StringArray `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

type ActivateUser struct {
	UserID      string    `db:"user_id"`
	IsActive    bool      `db:"is_active"`
	DateUpdated time.Time `db:"date_updated"`
}

type NewUser struct {
	Name     string
	Email    string
	Password string
}

type UpdateUser struct {
	ID       string
	Name     *string
	Email    *string
	Password string
}

type DeleteUser struct {
	ID       string
	Password string
}

type VerifyUser struct {
	Email string
	Token string
}

type VerifyToken struct {
	TokenID   string    `db:"token_id"`
	UserID    string    `db:"user_id"`
	Email     string    `db:"email"`
	ExpiresAt time.Time `db:"expires_at"`
	IssuedAt  time.Time `db:"issued_at"`
}
