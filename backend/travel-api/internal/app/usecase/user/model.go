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
	TokenVersion int32          `db:"token_version"`
	Roles        pq.StringArray `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

type NewUser struct {
	Name     string
	Email    string
	Password string
}

type UpdateUser struct {
	Name     *string
	Email    *string
	Password string
}
