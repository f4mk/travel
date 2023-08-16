package auth

import "errors"

var (
	ErrMissingKey           = errors.New("missing active key")
	ErrMissingKID           = errors.New("missing key id in header")
	ErrKIDFormat            = errors.New("key id must be string")
	ErrPrivateNotFound      = errors.New("missing private key for id")
	ErrSigningToken         = errors.New("error signing token")
	ErrLoadRevokedTokens    = errors.New("error loading revoked tokens")
	ErrParseToken           = errors.New("error parsing token")
	ErrValidateToken        = errors.New("error validating token")
	ErrExpiredToken         = errors.New("error expired token")
	ErrInvalidToken         = errors.New("error invalid token")
	ErrInvalidRefreshToken  = errors.New("error invalid refresh token")
	ErrCheckCachedToken     = errors.New("error checking token in cache")
	ErrRevokedToken         = errors.New("error revoked token")
	ErrEncodeTokenForCache  = errors.New("error encoding token for cache")
	ErrEncodeTokensForCache = errors.New("error encoding tokens for cache")
	ErrStoreCacheToken      = errors.New("error storing token in cache")
	ErrStoreCacheTokens     = errors.New("error storing tokens in cache")
	ErrReadTokensFromDB     = errors.New("error loading tokens from db")
	ErrReadTokenFromDB      = errors.New("error loading token from db")
	// middleware
	ErrAuthHeaderFormat = errors.New("error auth header format")
)
