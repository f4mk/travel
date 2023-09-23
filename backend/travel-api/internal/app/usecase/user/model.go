package user

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID           string
	Name         string
	Email        string
	IsActive     bool
	IsDeleted    bool
	TokenVersion int32
	Roles        pq.StringArray
	PasswordHash []byte
	DateCreated  time.Time
	DateUpdated  time.Time
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
	TokenID   string
	UserID    string
	Email     string
	ExpiresAt time.Time
	IssuedAt  time.Time
}
