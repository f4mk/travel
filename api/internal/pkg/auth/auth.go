package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
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
	AuthDuration    time.Duration
	RefreshDuration time.Duration
}

type Config struct {
	ActiveKIDs      []string
	KeyLookup       KeyLookup
	Cache           *redis.Client
	DB              *sqlx.DB
	AuthDuration    time.Duration
	RefreshDuration time.Duration
}

func New(cfg Config) (*Auth, error) {

	for _, activeKID := range cfg.ActiveKIDs {
		_, err := cfg.KeyLookup.PrivateKey(activeKID)
		if err != nil {
			return nil, fmt.Errorf("cannot find active key: %w", err)
		}
	}

	method := jwt.SigningMethodRS256

	keyFunc := func(t *jwt.Token) (any, error) {
		kid, ok := t.Header["kid"]

		if !ok {
			return nil, fmt.Errorf("missing key ID in token header")
		}

		kidStr, ok := kid.(string)

		if !ok {
			return nil, fmt.Errorf("key ID must be of type string")
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
		AuthDuration:    cfg.AuthDuration,
		RefreshDuration: cfg.RefreshDuration,
	}

	if err := a.LoadRevokedTokensToCache(); err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *Auth) GenerateToken(claims Claims, duration time.Duration) (string, error) {
	ia := time.Now()
	ea := ia.Add(duration)
	jti := uuid.New().String()
	claims.ID = jti
	claims.IssuedAt = &jwt.NumericDate{Time: ia}
	claims.ExpiresAt = &jwt.NumericDate{Time: ea}
	token := jwt.NewWithClaims(a.method, claims)
	// TODO: select the most fresh KID
	currentKID := a.activeKIDs[0]
	token.Header["kid"] = currentKID
	privateKey, err := a.keyLookup.PrivateKey(currentKID)
	if err != nil {
		return "", fmt.Errorf("no private key for key ID: %v", currentKID)
	}

	str, err := token.SignedString(privateKey)

	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return str, nil
}

func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &claims, a.keyFunc)

	if err != nil {
		return Claims{}, fmt.Errorf("cannot parse token: %w", err)
	}

	if !token.Valid {
		return Claims{}, errors.New("token is invalid")
	}

	return claims, nil
}

// TODO: refactor this to check is token was revoked
func (a *Auth) ValidateRefreshToken(ctx context.Context, tokenStr string) (Claims, error) {
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &claims, a.keyFunc)

	if err != nil {
		return Claims{}, fmt.Errorf("cannot parse token: %w", err)
	}

	if !token.Valid {
		return Claims{}, errors.New("token is invalid")
	}

	revoked, err := a.cache.Exists(ctx, claims.ID).Result()
	if err != nil {
		return Claims{}, fmt.Errorf("error checking token in cache: %w", err)
	}
	if revoked == 1 {
		return Claims{}, errors.New("token has been revoked")
	}

	return claims, nil
}

func (a *Auth) LoadRevokedTokensToCache() error {
	// Step 1: Query Database for Revoked Tokens
	revokedTokens, err := a.getRevokedTokens()
	if err != nil {
		return fmt.Errorf("error fetching revoked tokens from database: %w", err)
	}

	// Step 2: Insert Revoked Tokens into Redis
	for _, token := range revokedTokens {
		err := a.cache.Set(context.TODO(), token.ID, "revoked", time.Duration(token.ExpiresAt))
		if err != nil {
			return fmt.Errorf("error setting token in Redis: %w", err.Err())
		}
	}

	return nil
}

type Token struct {
	ID        string `db:"token_id"`
	ExpiresAt int64  `db:"expires_at"`
	IssuedAt  int64  `db:"issued_at"`
	RevokedAt int64  `db:"revoked_at"`
	Subject   string `db:"subject"`
}

func (a *Auth) getRevokedTokens() ([]Token, error) {
	var tokens []Token
	err := a.db.Select(&tokens, "SELECT token_id, expires_at, issued_at, revoked_at, subject FROM revoked_tokens")
	if err != nil {
		return nil, fmt.Errorf("error fetching revoked tokens from database: %w", err)
	}
	return tokens, nil
}
