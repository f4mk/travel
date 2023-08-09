package auth

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type Storer interface {
	// TODO: add CRUD methods
	Update(ctx context.Context, u User) error
	Delete(ctx context.Context, t string) error
	DeleteAll(ctx context.Context, uID string) error
	QueryByEmail(ctx context.Context, email string) (User, error)
}

type Core struct {
	storer Storer
	log    *zerolog.Logger
}

func NewCore(s Storer, l *zerolog.Logger) *Core {

	return &Core{
		storer: s,
		log:    l,
	}
}

func (c *Core) Login(ctx context.Context, email string, pw string, t string, exp time.Time) error {
	// TODO: check user exists, check pw, store token in tokens table with user id and exp date
	return nil
}

func (c *Core) Logout(ctx context.Context, email string, t string) error {
	// TODO: delete this token from tokens table
	return nil
}

func (c *Core) LogoutAll(ctx context.Context, email string) error {
	// TODO: delete all tokens from tokens table for that user
	return nil
}

func (c *Core) RevokeToken(ctx context.Context, email string, t string) error {
	// TODO: delete this token from tokens table
	return nil
}

func (c *Core) RevokeTokens(ctx context.Context, email string) error {
	// TODO: delete all tokens from tokens table for that user
	return nil
}

func (c *Core) RefreshToken(ctx context.Context, email string, t string) error {
	// TODO: find token in tokes, update the token with a new one
	return nil
}

func (c *Core) ResetPassword(ctx context.Context, email string) error {
	// TODO: send email with reset link
	return nil
}
