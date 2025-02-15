package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, user User) error
	Verify(ctx context.Context, user User) error
	QueryAll(ctx context.Context) ([]User, error)
	QueryByID(ctx context.Context, userID string) (User, error)
	QueryByEmail(ctx context.Context, email string) (User, error)
	QueryTokenByEmail(ctx context.Context, email string) (VerifyToken, error)
	DeleteVerifyTokensByUserID(ctx context.Context, userID string) error
	StoreVerifyToken(ctx context.Context, vt VerifyToken) error
}

// Core unit implements a set of methods for model types transformation.
// Core should neither be aware of a database implementation
// nor of the particular way of retrieving necessary data.
// Core dictates the interface that Storer dependency type must implement

type Core struct {
	storer Storer
	log    *zerolog.Logger
}

func NewCore(l *zerolog.Logger, s Storer) *Core {
	return &Core{
		storer: s,
		log:    l,
	}
}

func (c *Core) QueryAll(ctx context.Context) ([]User, error) {
	ctx, span := web.AddSpan(ctx, "usecase.user.query-all")
	defer span.End()
	tID := web.GetTraceID(ctx)
	us, err := c.storer.QueryAll(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: query all: %s", database.ErrQueryDB.Error())
		return []User{}, database.WrapStorerError(err)
	}
	return us, nil
}

func (c *Core) Create(ctx context.Context, nu NewUser) (User, string, error) {
	ctx, span := web.AddSpan(ctx, "usecase.user.create")
	defer span.End()
	tID := web.GetTraceID(ctx)
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: create: %s", auth.ErrGenHash.Error())
		return User{}, "", auth.ErrGenHash
	}
	now := time.Now().UTC()
	u := User{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		IsActive:     false,
		IsDeleted:    false,
		TokenVersion: 0,
		PasswordHash: hash,
		// TODO: may be find a better place
		Roles:       []string{auth.RoleUser},
		DateCreated: now,
		DateUpdated: now,
	}
	if err := c.storer.Create(ctx, u); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: create: %s", database.ErrQueryDB.Error())
		return User{}, "", database.WrapStorerError(err)
	}
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: create: %s", auth.ErrGenResetToken.Error())
		return User{}, "", auth.ErrGenResetToken
	}
	et := hex.EncodeToString(token)
	vt := VerifyToken{
		TokenID:   et,
		UserID:    u.ID,
		Email:     u.Email,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		IssuedAt:  time.Now().UTC(),
	}
	// TODO: make a cron to clear the table
	err = c.storer.StoreVerifyToken(ctx, vt)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: create: %s", database.ErrQueryDB.Error())
		return User{}, "", database.WrapStorerError(err)
	}
	return u, et, nil
}

func (c *Core) Update(ctx context.Context, uu UpdateUser) (User, error) {
	ctx, span := web.AddSpan(ctx, "usecase.user.update")
	defer span.End()
	tID := web.GetTraceID(ctx)
	u, err := c.storer.QueryByID(ctx, uu.ID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: update: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(uu.Password)); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: update: %s", web.ErrAuthFailed.Error())
		return User{}, web.ErrAuthFailed
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: update: %s", auth.ErrGetClaims.Error())
		return User{}, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && uu.ID != u.ID {
		c.log.Error().Str("TraceID", tID).Msgf("user: update: %s", web.ErrForbidden.Error())
		return User{}, web.ErrForbidden
	}
	if uu.Name != nil {
		u.Name = *uu.Name
	}
	if uu.Email != nil {
		u.Email = *uu.Email
	}
	u.DateUpdated = time.Now().UTC()
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: update: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	return u, nil
}

func (c *Core) Verify(ctx context.Context, vu VerifyUser) (User, error) {
	ctx, span := web.AddSpan(ctx, "usecase.user.verify")
	defer span.End()
	tID := web.GetTraceID(ctx)
	vt, err := c.storer.QueryTokenByEmail(ctx, vu.Email)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: verify: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	if vt.ExpiresAt.Before(time.Now().UTC()) {
		c.log.Error().Str("TraceID", tID).Msgf("user: verify: %s", auth.ErrValidateVerifyToken.Error())
		return User{}, auth.ErrValidateResetToken
	}
	// delete all tokens
	if err := c.storer.DeleteVerifyTokensByUserID(ctx, vt.UserID); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: verify: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	u, err := c.storer.QueryByID(ctx, vt.UserID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: verify: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	u.IsActive = true
	u.DateUpdated = time.Now().UTC()
	if err := c.storer.Verify(ctx, u); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: verify: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	return u, nil
}

func (c *Core) QueryByID(ctx context.Context, uID string) (User, error) {
	ctx, span := web.AddSpan(ctx, "usecase.user.querry-by-id")
	defer span.End()
	tID := web.GetTraceID(ctx)
	u, err := c.storer.QueryByID(ctx, uID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: query by id: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: query by id: %s", auth.ErrGetClaims.Error())
		return User{}, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && uID != u.ID {
		c.log.Error().Str("TraceID", tID).Msgf("user: query by id: %s", web.ErrForbidden.Error())
		return User{}, web.ErrForbidden
	}
	return u, nil
}

func (c *Core) Delete(ctx context.Context, du DeleteUser) (User, error) {
	ctx, span := web.AddSpan(ctx, "usecase.user.delete")
	defer span.End()
	tID := web.GetTraceID(ctx)
	u, err := c.storer.QueryByID(ctx, du.ID)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: delete: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(du.Password)); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("auth: login: %s", web.ErrAuthFailed.Error())
		return User{}, web.ErrAuthFailed
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: delete: %s", auth.ErrGetClaims.Error())
		return User{}, auth.ErrGetClaims
	}
	if !claims.Authorize(auth.RoleAdmin) && claims.Subject != du.ID {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: delete: %s", web.ErrForbidden.Error())
		return User{}, web.ErrForbidden
	}
	u.DateUpdated = time.Now().UTC()
	u.IsActive = false
	u.IsDeleted = true
	u.TokenVersion = u.TokenVersion + 1
	if err := c.storer.Delete(ctx, u); err != nil {
		c.log.Err(err).Str("TraceID", tID).Msgf("user: delete: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	return u, nil
}
