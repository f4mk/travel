package user

import (
	"time"

	"github.com/lib/pq"
)

// Model describes a data type in terms of business needs.
type User struct {
	ID           string         `db:"user_id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Email        string         `db:"email" json:"email"`
	Roles        pq.StringArray `db:"roles" json:"-"`
	PasswordHash []byte         `db:"password_hash" json:"-"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
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
