package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// lib/pq errorCodeNames
// https://github.com/lib/pq/blob/master/error.go#L178
const (
	uniqueViolation = pq.ErrorCode("23505")
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrForbidden     = errors.New("not allowed")
	ErrAuthFailed    = errors.New("authentication failed")
	ErrAlreadyExists = errors.New("already exists")
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

	usr := User{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := c.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", wrapError(err))
	}

	return usr, nil
}

func (c *Core) Update(ctx context.Context, uID string, uu UpdateUser) (User, error) {

	//query existing user
	usr, err := c.storer.QueryByID(ctx, uID)

	if err != nil {
		return User{}, fmt.Errorf("query user: %w", wrapError(err))
	}

	claims, err := auth.GetClaims(ctx)

	if err != nil {
		return User{}, fmt.Errorf("get claims: %w", err)
	}

	// TODO: should check ID in JWT token
	if !claims.Authorize(auth.RoleAdmin) && uID != usr.ID {
		return User{}, ErrForbidden
	}

	//update user
	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generate password hash: %w", err)
		}
		usr.PasswordHash = hash
	}
	if len(uu.Roles) > 0 {
		usr.Roles = uu.Roles
	}
	usr.DateUpdated = time.Now().UTC()

	if err := c.storer.Update(ctx, usr); err != nil {
		return User{}, fmt.Errorf("update: %w", err)
	}
	return usr, nil
}

func (c *Core) QueryByID(ctx context.Context, uID string) (User, error) {

	usr, err := c.storer.QueryByID(ctx, uID)
	if err != nil {
		return User{}, fmt.Errorf("query user: %w", wrapError(err))
	}

	return usr, nil
}

func (c *Core) Delete(ctx context.Context, uID string) error {

	//query existing user
	_, err := c.storer.QueryByID(ctx, uID)
	if err != nil {
		return fmt.Errorf("query user: %w", wrapError(err))
	}

	claims, err := auth.GetClaims(ctx)

	if err != nil {
		return fmt.Errorf("get claims: %w", err)
	}

	// TODO: should check ID in JWT token == uID
	if !claims.Authorize(auth.RoleAdmin) {
		return ErrForbidden
	}

	if err := c.storer.Delete(ctx, uID); err != nil {
		return fmt.Errorf("delete user: %w", wrapError(err))
	}
	return nil
}

// TODO: may be should be here
func wrapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolation {
		return ErrAlreadyExists
	}

	return err
}
