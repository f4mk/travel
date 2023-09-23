package auth

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
	TokenID      string
	Subject      string
	TokenVersion int32
	IssuedAt     time.Time
	ExpiresAt    time.Time
	RevokedAt    time.Time
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
	TokenID   string
	UserID    string
	Email     string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

type SubmitPassword struct {
	ResetToken string
	Password   string
}
