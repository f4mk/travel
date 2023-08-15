package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/f4mk/api/internal/pkg/database"
	"github.com/f4mk/api/pkg/web"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	// TODO: add CRUD methods
	DeleteToken(ctx context.Context, dt DeleteToken) error
	DeleteAllTokes(ctx context.Context, uID string) error
	QueryByEmail(ctx context.Context, email string) (User, error)
	QueryByID(ctx context.Context, uID string) (User, error)
	Update(ctx context.Context, u User) error
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

func (c *Core) Login(ctx context.Context, lu LoginUser) (AuthenticatedUser, error) {
	u, err := c.storer.QueryByEmail(ctx, lu.Email)
	if err != nil {
		// return ErrAuthFailed to not spoil user email if not found
		return AuthenticatedUser{}, fmt.Errorf("wrong credentials: %w", web.ErrAuthFailed)
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(lu.Password)); err != nil {
		return AuthenticatedUser{}, fmt.Errorf("wrong credentials: %w", web.ErrAuthFailed)
	}
	au := AuthenticatedUser{
		UserID: u.UserID,
		Email:  u.Email,
		Name:   u.Name,
		Roles:  u.Roles,
	}
	return au, nil
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

func (c *Core) ChangePassword(ctx context.Context, cp ChangePassword) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cp.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generate password hash: %w", err)
	}
	u, err := c.storer.QueryByID(ctx, cp.UserID)
	if err != nil {
		return User{}, fmt.Errorf("query user: %w", database.WrapBusinessError(err))
	}
	u.PasswordHash = hash
	u.DateUpdated = time.Now().UTC()
	if err := c.storer.Update(ctx, u); err != nil {
		return User{}, fmt.Errorf("update: %w", database.WrapBusinessError(err))
	}
	if err := c.storer.DeleteAllTokes(ctx, cp.UserID); err != nil {
		return User{}, fmt.Errorf("update: %w", database.WrapBusinessError(err))
	}
	return u, nil
}

//revive:disable
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

//revive:enable
