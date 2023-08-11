package user

import (
	"time"

	"github.com/f4mk/api/pkg/web"
	"github.com/lib/pq"
)

// Model describes a data type in terms of business needs.
type User struct {
	// TODO: validation
	ID           string         `db:"user_id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Email        string         `db:"email" json:"email"`
	Roles        pq.StringArray `db:"roles" json:"-"`
	PasswordHash []byte         `db:"password_hash" json:"-"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
}

type NewUser struct {
	// TODO: validation
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

func (nu NewUser) Validate() error {
	if err := web.Check(nu); err != nil {
		return err
	}
	return nil
}

type UpdateUser struct {
	// TODO: validation
	Name            *string `json:"name"`
	Email           *string `json:"email" validate:"omitempty,email"`
	Password        *string `json:"password" validate:"omitempty"`
	PasswordConfirm *string `json:"password_confirm" validate:"eqfield=Password"`
}

func (uu UpdateUser) Validate() error {
	if err := web.Check(uu); err != nil {
		return err
	}
	return nil
}
