package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/f4mk/travel/backend/travel-api/internal/pkg/auth"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/database"
	"github.com/f4mk/travel/backend/travel-api/internal/pkg/web"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	DeleteToken(ctx context.Context, dt DeleteToken) error
	StoreResetToken(ctx context.Context, rt ResetToken) error
	DeleteResetTokensByUserID(ctx context.Context, uID string) error
	QueryResetTokenByID(ctx context.Context, token string) (ResetToken, error)
	QueryByEmail(ctx context.Context, email string) (User, error)
	QueryByID(ctx context.Context, uID string) (User, error)
	Update(ctx context.Context, u User) error
}

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

func (c *Core) Login(ctx context.Context, lu LoginUser) (AuthenticatedUser, error) {
	u, err := c.storer.QueryByEmail(ctx, lu.Email)
	if err != nil {
		// NOTE: return ErrAuthFailed to not spoil user email if not found
		c.log.Err(err).Msgf("auth: login: %s", database.ErrQueryDB.Error())
		return AuthenticatedUser{}, web.ErrAuthFailed
	}
	if !u.IsActive {
		c.log.Error().Msgf("auth: login: user is inactive")
		return AuthenticatedUser{}, web.ErrAuthFailed
	}
	if !u.IsDeleted {
		c.log.Error().Msgf("auth: login: user is deleted")
		return AuthenticatedUser{}, web.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(lu.Password)); err != nil {
		c.log.Err(err).Msgf("auth: login: %s", web.ErrAuthFailed.Error())
		return AuthenticatedUser{}, web.ErrAuthFailed
	}
	au := AuthenticatedUser{
		UserID:       u.UserID,
		Email:        u.Email,
		Name:         u.Name,
		TokenVersion: u.TokenVersion,
		Roles:        u.Roles,
		DateCreated:  u.DateCreated,
	}
	return au, nil
}

func (c *Core) Logout(ctx context.Context, dt DeleteToken) error {
	_, err := c.storer.QueryByID(ctx, dt.Subject)
	if err != nil {
		c.log.Err(err).Msgf("auth: logout: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	if err := c.storer.DeleteToken(ctx, dt); err != nil {
		c.log.Err(err).Msgf("auth: logout: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}

	return nil
}

func (c *Core) ChangePassword(ctx context.Context, cp ChangePassword) (User, error) {
	u, err := c.storer.QueryByID(ctx, cp.UserID)
	if err != nil {
		c.log.Err(err).Msgf("auth: change password: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(cp.PasswordOld)); err != nil {
		c.log.Err(err).Msgf("auth: login: %s", web.ErrAuthFailed.Error())
		return User{}, web.ErrAuthFailed
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(cp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.Err(err).Msgf("auth: change password: %s", auth.ErrGenHash.Error())
		return User{}, auth.ErrGenHash
	}
	// update user with new password and token version
	u.PasswordHash = hash
	u.DateUpdated = time.Now().UTC()
	u.TokenVersion = u.TokenVersion + 1
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Msgf("auth: change password: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	return u, nil
}

// TODO: cron to clean the table every day for unclaimed tokens
func (c *Core) ResetPasswordRequest(ctx context.Context, email string) (ResetPassword, error) {
	u, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password request: %s", database.ErrQueryDB.Error())
		return ResetPassword{}, database.WrapStorerError(err)
	}

	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password request: %s", auth.ErrGenResetToken.Error())
		return ResetPassword{}, auth.ErrGenResetToken
	}
	et := hex.EncodeToString(token)
	rt := ResetToken{
		TokenID:   et,
		UserID:    u.UserID,
		Email:     u.Email,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		IssuedAt:  time.Now().UTC(),
	}
	err = c.storer.StoreResetToken(ctx, rt)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password request: %s", database.ErrQueryDB.Error())
		return ResetPassword{}, database.WrapStorerError(err)
	}
	rp := ResetPassword{
		Email:      u.Email,
		Name:       u.Name,
		ResetToken: et,
	}
	return rp, nil
}

func (c *Core) ResetPasswordSubmit(ctx context.Context, sp SubmitPassword) (User, error) {
	rt, err := c.storer.QueryResetTokenByID(ctx, sp.ResetToken)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password validate: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	if rt.ExpiresAt.Before(time.Now().UTC()) {
		c.log.Error().Msgf("auth: reset password validate: %s", auth.ErrValidateResetToken.Error())
		return User{}, auth.ErrValidateResetToken
	}
	// delete all tokens
	if err := c.storer.DeleteResetTokensByUserID(ctx, rt.UserID); err != nil {
		c.log.Err(err).Msgf("auth: reset password validate: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	u, err := c.storer.QueryByEmail(ctx, rt.Email)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password validate: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	if u.IsDeleted {
		return User{}, web.ErrNotFound
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(sp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password submit: %s", auth.ErrGenHash.Error())
		return User{}, auth.ErrGenHash
	}
	// update user with new password and token version
	u.PasswordHash = hash
	u.DateUpdated = time.Now().UTC()
	u.TokenVersion = u.TokenVersion + 1
	// NOTE: if user didnt get verification email, they can ask for reset pwd
	// TODO: maybe need a separate handler
	u.IsActive = true
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Msgf("auth: reset password submit: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	return u, nil
}

func (c *Core) LogoutAll(ctx context.Context, dt DeleteToken) (int32, error) {
	u, err := c.storer.QueryByID(ctx, dt.Subject)
	if err != nil {
		c.log.Err(err).Msgf("auth: logout all: %s", database.ErrQueryDB.Error())
		return 0, database.WrapStorerError(err)
	}
	u.DateUpdated = time.Now().UTC()
	u.TokenVersion = u.TokenVersion + 1
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Msgf("auth: logout all: %s", database.ErrQueryDB.Error())
		return 0, database.WrapStorerError(err)
	}
	return u.TokenVersion, nil
}
