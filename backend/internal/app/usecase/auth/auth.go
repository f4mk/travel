package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/f4mk/api/internal/pkg/auth"
	"github.com/f4mk/api/internal/pkg/database"
	"github.com/f4mk/api/pkg/web"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	DeleteToken(ctx context.Context, dt DeleteToken) error
	StoreResetToken(ctx context.Context, rt ResetToken) error
	UpdateResetToken(ctx context.Context, rt ResetToken) error
	RetrieveResetToken(ctx context.Context, token string) (ResetToken, error)
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
		// NOTE: return ErrAuthFailed to not spoil user email if not found
		c.log.Err(err).Msgf("auth: login: %s", database.ErrQueryDB.Error())
		return AuthenticatedUser{}, web.ErrAuthFailed
	}
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(lu.Password)); err != nil {
		c.log.Err(err).Msgf("auth: login: %s", web.ErrAuthFailed.Error())
		return AuthenticatedUser{}, web.ErrAuthFailed
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
	hash, err := bcrypt.GenerateFromPassword([]byte(cp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.Err(err).Msgf("auth: change password: %s", auth.ErrGenHash.Error())
		return User{}, auth.ErrGenHash
	}
	u, err := c.storer.QueryByID(ctx, cp.UserID)
	if err != nil {
		c.log.Err(err).Msgf("auth: change password: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	u.PasswordHash = hash
	u.DateUpdated = time.Now().UTC()
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Msgf("auth: change password: %s", database.ErrQueryDB.Error())
		return User{}, database.WrapStorerError(err)
	}
	if err := c.storer.DeleteAllTokes(ctx, cp.UserID); err != nil {
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
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password request: %s", auth.ErrGenResetToken.Error())
		return ResetPassword{}, auth.ErrGenResetToken
	}
	et := hex.EncodeToString(token)
	rt := ResetToken{
		TokenID:   et,
		Email:     u.Email,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		IssuedAt:  time.Now().UTC(),
		Used:      false,
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

func (c *Core) ResetPasswordSubmit(ctx context.Context, sp SubmitPassword) error {
	rt, err := c.storer.RetrieveResetToken(ctx, sp.ResetToken)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password validate: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	if rt.ExpiresAt.Before(time.Now().UTC()) {
		c.log.Error().Msgf("auth: reset password validate: %s", auth.ErrValidateResetToken.Error())
		return auth.ErrValidateResetToken
	}
	// mark token as used
	rt.Used = true
	if err := c.storer.UpdateResetToken(ctx, rt); err != nil {
		c.log.Err(err).Msgf("auth: reset password validate: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	u, err := c.storer.QueryByEmail(ctx, rt.Email)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password validate: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(sp.Password), bcrypt.DefaultCost)
	if err != nil {
		c.log.Err(err).Msgf("auth: reset password submit: %s", auth.ErrGenHash.Error())
		return auth.ErrGenHash
	}
	// update user with new password
	u.PasswordHash = hash
	u.DateUpdated = time.Now().UTC()
	if err := c.storer.Update(ctx, u); err != nil {
		c.log.Err(err).Msgf("auth: reset password submit: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	if err := c.storer.DeleteAllTokes(ctx, u.UserID); err != nil {
		c.log.Err(err).Msgf("auth: reset password submit: %s", database.ErrQueryDB.Error())
		return database.WrapStorerError(err)
	}
	return nil
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
