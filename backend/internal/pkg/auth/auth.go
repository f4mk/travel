package auth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type KeyLookup interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	PublicKey(kid string) (*rsa.PublicKey, error)
}

type Auth struct {
	activeKIDs      []string
	method          jwt.SigningMethod
	keyFunc         func(t *jwt.Token) (any, error)
	parser          *jwt.Parser
	keyLookup       KeyLookup
	cache           *redis.Client
	db              *sqlx.DB
	log             *zerolog.Logger
	AuthDuration    time.Duration
	RefreshDuration time.Duration
}

type Config struct {
	ActiveKIDs      []string
	KeyLookup       KeyLookup
	Cache           *redis.Client
	DB              *sqlx.DB
	Log             *zerolog.Logger
	AuthDuration    time.Duration
	RefreshDuration time.Duration
}

func New(cfg Config) (*Auth, error) {

	for _, activeKID := range cfg.ActiveKIDs {
		_, err := cfg.KeyLookup.PrivateKey(activeKID)
		if err != nil {
			cfg.Log.Err(err).Msg(ErrMissingKey.Error())
			return nil, ErrMissingKey
		}
	}

	method := jwt.SigningMethodRS256

	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]

		if !ok {
			cfg.Log.Error().Msg(ErrMissingKID.Error())
			return nil, ErrMissingKID
		}

		kidStr, ok := kid.(string)

		if !ok {
			cfg.Log.Error().Msg(ErrKIDFormat.Error())
			return nil, ErrKIDFormat
		}

		return cfg.KeyLookup.PublicKey(kidStr)
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{method.Name}))

	a := Auth{
		method:          jwt.GetSigningMethod(method.Name),
		keyLookup:       cfg.KeyLookup,
		keyFunc:         keyFunc,
		parser:          parser,
		activeKIDs:      cfg.ActiveKIDs,
		cache:           cfg.Cache,
		db:              cfg.DB,
		log:             cfg.Log,
		AuthDuration:    cfg.AuthDuration,
		RefreshDuration: cfg.RefreshDuration,
	}

	if err := a.LoadRevokedTokensToCache(); err != nil {
		a.log.Err(err).Msg(ErrLoadRevokedTokens.Error())
		return nil, ErrLoadRevokedTokens
	}

	return &a, nil
}

func (a *Auth) GenerateToken(ctx context.Context, claims Claims, duration time.Duration) (string, error) {
	ia := time.Now().UTC()
	ea := ia.Add(duration)
	jti := uuid.New().String()
	claims.ID = jti
	tv, err := a.getLastTokenVersion(ctx, claims)
	if err != nil {
		a.log.Err(err).Msg(ErrGetClaims.Error())
		return "", ErrGetClaims
	}
	claims.TokenVersion = tv
	claims.IssuedAt = &jwt.NumericDate{Time: ia}
	claims.ExpiresAt = &jwt.NumericDate{Time: ea}
	token := jwt.NewWithClaims(a.method, claims)
	// TODO: select the most fresh KID
	currentKID := a.activeKIDs[0]
	token.Header["kid"] = currentKID
	privateKey, err := a.keyLookup.PrivateKey(currentKID)
	if err != nil {
		a.log.Err(err).Msg(ErrPrivateNotFound.Error())
		return "", ErrPrivateNotFound
	}

	str, err := token.SignedString(privateKey)

	if err != nil {
		a.log.Err(err).Msg(ErrSigningToken.Error())
		return "", ErrSigningToken
	}

	return str, nil
}

func (a *Auth) ValidateToken(ctx context.Context, t string) (Claims, error) {
	var claims Claims
	token, err := a.parser.ParseWithClaims(t, &claims, a.keyFunc)

	if err != nil {
		a.log.Err(err).Msg(ErrParseToken.Error())
		return Claims{}, ErrParseToken
	}
	revoked, err := a.cache.Exists(ctx, claims.ID).Result()
	if err != nil {
		a.log.Err(err).Msg(ErrCheckCachedToken.Error())
		return Claims{}, ErrCheckCachedToken
	}
	if revoked == 1 {
		a.log.Error().Msg(ErrRevokedToken.Error())
		return Claims{}, ErrRevokedToken
	}
	et, err := token.Claims.GetExpirationTime()
	if err != nil {
		a.log.Err(err).Msg(ErrValidateToken.Error())
		return Claims{}, ErrValidateToken
	}
	if et.Time.Before(time.Now()) {
		a.log.Error().Msg(ErrExpiredToken.Error())
		return Claims{}, ErrExpiredToken
	}
	tv, err := a.getLastTokenVersion(ctx, claims)
	if err != nil {
		a.log.Err(err).Msg(ErrValidateToken.Error())
		return Claims{}, ErrValidateToken
	}
	if claims.TokenVersion != tv {
		a.log.Error().Msg(ErrValidateTokenVersion.Error())
		return Claims{}, ErrValidateTokenVersion
	}
	if !token.Valid {
		a.log.Error().Msg(ErrInvalidToken.Error())
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}

func (a *Auth) StoreUserTokenVersion(ctx context.Context, uID string, tv int32) error {
	if err := a.cache.Set(
		ctx,
		uID,
		strconv.Itoa(int(tv)),
		0,
	).Err(); err != nil {
		a.log.Err(err).Msg(ErrStoreCacheTokenVersion.Error())
		return ErrStoreCacheTokenVersion
	}
	return nil
}

func (a *Auth) MarkTokenAsRevoked(ctx context.Context, t TokenParams) error {
	jsonData, err := json.Marshal(t)
	if err != nil {
		a.log.Err(err).Msg(ErrEncodeTokenForCache.Error())
		return ErrEncodeTokenForCache
	}
	if err := a.cache.Set(
		ctx,
		t.TokenID,
		jsonData,
		t.ExpiresAt.Sub(time.Now().UTC()),
	).Err(); err != nil {
		a.log.Err(err).Msg(ErrStoreCacheToken.Error())
		return ErrStoreCacheToken
	}
	return nil
}

func (a *Auth) LoadRevokedTokensToCache() error {
	// Step 1: Query Database for Revoked Tokens
	revokedTokens, err := a.getRevokedTokens()
	if err != nil {
		a.log.Err(err).Msg(ErrReadTokensFromDB.Error())
		return ErrReadTokensFromDB
	}

	// Step 2: Insert Revoked Tokens into Redis
	for _, t := range revokedTokens {
		jsonData, err := json.Marshal(t)
		if err != nil {
			a.log.Err(err).Msg(ErrEncodeTokensForCache.Error())
			return ErrEncodeTokensForCache
		}

		duration := t.ExpiresAt.Sub(time.Now().UTC())
		if duration < 0 {
			continue
		}

		if err := a.cache.Set(
			context.TODO(),
			t.TokenID,
			jsonData,
			duration).Err(); err != nil {
			a.log.Err(err).Msg(ErrStoreCacheTokens.Error())
			return ErrStoreCacheTokens
		}
	}

	return nil
}

type TokenParams struct {
	TokenID      string    `db:"token_id" json:"token_id"`
	Subject      string    `db:"subject" json:"subject"`
	TokenVersion int32     `db:"token_version" json:"token_version"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
	IssuedAt     time.Time `db:"issued_at" json:"issued_at"`
	RevokedAt    time.Time `db:"revoked_at" json:"revoked_at"`
}

func (a *Auth) getRevokedTokens() ([]TokenParams, error) {
	var tokens []TokenParams
	err := a.db.Select(&tokens, "SELECT token_id, subject, token_version, expires_at, issued_at, revoked_at  FROM revoked_tokens")
	if err != nil {
		a.log.Err(err).Msg(ErrReadTokenFromDB.Error())
		return nil, ErrReadTokenFromDB
	}
	return tokens, nil
}

type user struct {
	TokenVersion int32 `db:"token_version"`
}

func (a *Auth) getLastTokenVersion(ctx context.Context, c Claims) (int32, error) {

	res, err := a.cache.Get(ctx, c.Subject).Result()
	if err != nil {
		a.log.Warn().Msg(ErrCheckCachedTokenVersion.Error())

		u := user{}
		q := "SELECT token_version FROM users WHERE user_id = $1;"
		if err := a.db.Get(&u, q, c.Subject); err != nil {
			a.log.Err(err).Msg(ErrReadUserFromDB.Error())
			return 0, ErrReadUserFromDB
		}
		a.cache.Set(ctx, c.Subject, u.TokenVersion, 0)
		return u.TokenVersion, nil
	}

	valInt64, err := strconv.ParseInt(res, 10, 32)
	if err != nil {
		a.log.Err(err).Msg(ErrParseCachedTokenVersion.Error())
		return 0, ErrParseCachedTokenVersion
	}

	valInt32 := int32(valInt64)

	return valInt32, nil
}
