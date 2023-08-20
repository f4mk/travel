package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/database"
	"github.com/f4mk/api/pkg/web"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

const emailCooldown = 10

type Storer interface {
	StoreResetToken(ctx context.Context, rt ResetToken) error
	DeleteResetTokensByUserID(ctx context.Context, uID string) error
	QueryResetTokenByID(ctx context.Context, token string) (ResetToken, error)
	QueryLastResetTokenByUserID(ctx context.Context, uID string) (*ResetToken, error)
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
		// NOTE: return ErrAuthFailed to not spoil user email if not found
		c.log.Err(err).Msgf("auth: login: %s", database.ErrQueryDB.Error())
		return AuthenticatedUser{}, web.ErrAuthFailed
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
	}
	return au, nil
}

func (c *Core) Logout(ctx context.Context, dt DeleteToken) (int32, error) {
	u, err := c.storer.QueryByID(ctx, dt.Subject)
	if err != nil {
		c.log.Err(err).Msgf("auth: logout: %s", database.ErrQueryDB.Error())
		return 0, database.WrapStorerError(err)
	}
	u.DateUpdated = time.Now().UTC()
	u.TokenVersion = u.TokenVersion + 1
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Msgf("auth: logout: %s", database.ErrQueryDB.Error())
		return 0, database.WrapStorerError(err)
	}
	return u.TokenVersion, nil
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

func (c *Core) ResetPasswordRequest(ctx context.Context, email string) (ResetPassword, error) {
	u, err := c.storer.QueryByEmail(ctx, email)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password request: %s", database.ErrQueryDB.Error())
		return ResetPassword{}, database.WrapStorerError(err)
	}
	lt, err := c.storer.QueryLastResetTokenByUserID(ctx, u.UserID)
	if err != nil {
		//handle all errors but 404(no token - no problems)
		if !errors.Is(database.WrapStorerError(err), web.ErrNotFound) {
			c.log.Err(err).Msgf("auth: reset password request: %s", database.ErrQueryDB.Error())
			return ResetPassword{}, database.WrapStorerError(err)
		}
	}
	if lt != nil {
		// allow new token once every Xmin
		if lt.IssuedAt.After(time.Now().Add(-emailCooldown * time.Minute)) {
			c.log.Warn().Msg("auth: reset password request: requested token too soon")
			return ResetPassword{}, auth.ErrResetTokenReqLimit
		}
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
	hash, err := bcrypt.GenerateFromPassword([]byte(sp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password submit: %s", auth.ErrGenHash.Error())
		return User{}, auth.ErrGenHash
	}
	// update user with new password and token version
	u.PasswordHash = hash
	u.DateUpdated = time.Now().UTC()
	u.TokenVersion = u.TokenVersion + 1
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Msgf("auth: reset password submit: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
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

//revive:enable
