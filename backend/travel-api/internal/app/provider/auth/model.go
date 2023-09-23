package auth

import (
	"time"

	"github.com/lib/pq"
)

type StorerUser struct {
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
type StorerResetToken struct {
	TokenID   string    `db:"token_id"`
	UserID    string    `db:"user_id"`
	Email     string    `db:"email"`
	IssuedAt  time.Time `db:"issued_at"`
	ExpiresAt time.Time `db:"expires_at"`
}
type StorerDeleteToken struct {
	TokenID      string    `db:"token_id"`
	Subject      string    `db:"subject"`
	TokenVersion int32     `db:"token_version"`
	IssuedAt     time.Time `db:"issued_at"`
	ExpiresAt    time.Time `db:"expires_at"`
	RevokedAt    time.Time `db:"revoked_at"`
}
