package user

import (
	"context"
	"fmt"
	"time"

	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/database"
	"github.com/f4mk/api/pkg/web"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	// TODO:
	Create(ctx context.Context, u User) error
	Update(ctx context.Context, u User) error
	Delete(ctx context.Context, uID string) error
	QueryAll(ctx context.Context) ([]User, error)
	QueryByID(ctx context.Context, uID string) (User, error)
	// QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// Core unit implements a set of methods for model types transformation.
// Core should neither be aware of a database implementation
// nor of the particular way of retrieving necessary data.
// Core dictates the interface that Storer dependency type must implement

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

func (c *Core) QueryAll(ctx context.Context) ([]User, error) {
	return c.storer.QueryAll(ctx)
}

func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generate password hash: %w", err)
	}
	now := time.Now().UTC()
	u := User{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		// TODO: may be find a better place
		Roles:       []string{auth.RoleUser},
		DateCreated: now,
		DateUpdated: now,
	}
	if err := c.storer.Create(ctx, u); err != nil {
		return User{}, fmt.Errorf("create: %w", database.WrapBusinessError(err))
	}
	return u, nil
}

func (c *Core) Update(ctx context.Context, uID string, uu UpdateUser) (User, error) {
	u, err := c.storer.QueryByID(ctx, uID)
	if err != nil {
		return User{}, fmt.Errorf("query user: %w", database.WrapBusinessError(err))
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return User{}, fmt.Errorf("get claims: %w", err)
	}
	// TODO: should check ID in JWT token
	if !claims.Authorize(auth.RoleAdmin) && uID != u.ID {
		return User{}, web.ErrForbidden
	}
	//update user
	if uu.Name != nil {
		u.Name = *uu.Name
	}
	if uu.Email != nil {
		u.Email = *uu.Email
	}
	if uu.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generate password hash: %w", err)
		}
		u.PasswordHash = hash
	}
	u.DateUpdated = time.Now().UTC()
	if err := c.storer.Update(ctx, u); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}
	return u, nil
}

func (c *Core) QueryByID(ctx context.Context, uID string) (User, error) {
	u, err := c.storer.QueryByID(ctx, uID)
	if err != nil {
		return User{}, fmt.Errorf("query user: %w", database.WrapBusinessError(err))
	}
	return u, nil
}

func (c *Core) Delete(ctx context.Context, uID string) error {
	_, err := c.storer.QueryByID(ctx, uID)
	if err != nil {
		return fmt.Errorf("query user: %w", database.WrapBusinessError(err))
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return fmt.Errorf("get claims: %w", err)
	}
	if !claims.Authorize(auth.RoleAdmin) || claims.Subject != uID {
		return web.ErrForbidden
	}
	if err := c.storer.Delete(ctx, uID); err != nil {
		return fmt.Errorf("delete user: %w", database.WrapBusinessError(err))
	}
	return nil
}
