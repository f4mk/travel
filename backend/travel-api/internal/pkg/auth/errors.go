package auth

import "errors"

var (
	ErrMissingKey              = errors.New("missing active key")
	ErrMissingKID              = errors.New("missing key id in header")
	ErrKIDFormat               = errors.New("key id must be string")
	ErrPrivateNotFound         = errors.New("missing private key for id")
	ErrSigningToken            = errors.New("error signing token")
	ErrLoadRevokedTokens       = errors.New("error loading revoked tokens")
	ErrParseToken              = errors.New("error parsing token")
	ErrValidateToken           = errors.New("error validating token")
	ErrValidateTokenVersion    = errors.New("error validating token version")
	ErrExpiredToken            = errors.New("error expired token")
	ErrInvalidToken            = errors.New("error invalid token")
	ErrInvalidRefreshToken     = errors.New("error invalid refresh token")
	ErrCheckCachedToken        = errors.New("error checking token in cache")
	ErrCheckCachedTokenVersion = errors.New("error checking token version in cache")
	ErrParseCachedTokenVersion = errors.New("error parsing token version in cache")
	ErrRevokedToken            = errors.New("error revoked token")
	ErrEncodeTokenForCache     = errors.New("error encoding token for cache")
	ErrEncodeTokensForCache    = errors.New("error encoding tokens for cache")
	ErrStoreCacheToken         = errors.New("error storing token in cache")
	ErrStoreCacheTokenVersion  = errors.New("error storing token version in cache")
	ErrStoreCacheTokens        = errors.New("error storing tokens in cache")
	ErrReadTokensFromDB        = errors.New("error loading tokens from db")
	ErrReadTokenFromDB         = errors.New("error loading token from db")
	ErrReadUserFromDB          = errors.New("error loading user from db")
	ErrGenHash                 = errors.New("error generate hash")
	ErrGenResetToken           = errors.New("error generate reset token")
	ErrValidateResetToken      = errors.New("error validate reset token")
	ErrValidateVerifyToken     = errors.New("error validate verify token")
	ErrResetTokenReqLimit      = errors.New("error request reset token too often")
	ErrGetClaims               = errors.New("error get claims")
	ErrAuthHeaderFormat        = errors.New("error auth header format")
)
