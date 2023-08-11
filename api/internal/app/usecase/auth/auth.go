package auth

import (
	"context"
	"fmt"

	"github.com/f4mk/api/internal/pkg/database"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	// TODO: add CRUD methods
	DeleteToken(ctx context.Context, dt DeleteToken) error
	DeleteAllTokes(ctx context.Context, uID string) error
	QueryByEmail(ctx context.Context, email string) (User, error)
	QueryByID(ctx context.Context, uID string) (User, error)
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

func (c *Core) Login(ctx context.Context, u LoginUser) (AuthenticatedUser, error) {
	user, err := c.storer.QueryByEmail(ctx, u.Email)

	if err != nil {
		return AuthenticatedUser{}, fmt.Errorf("query user: %w", database.WrapBusinessError(err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(u.Password)); err != nil {
		return AuthenticatedUser{}, fmt.Errorf("wrong credentials: %w", database.WrapBusinessError(err))
	}
	auser := AuthenticatedUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Roles: user.Roles,
	}

	return auser, nil
}

func (c *Core) Logout(ctx context.Context, dt DeleteToken) error {
	_, err := c.storer.QueryByID(ctx, dt.Subject)
	if err != nil {
		return fmt.Errorf("query user: %w", database.WrapBusinessError(err))
	}
	if err := c.storer.DeleteToken(ctx, dt); err != nil {
		return fmt.Errorf("delete token: %w", database.WrapBusinessError(err))
	}

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
