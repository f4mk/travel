package auth

import "errors"

var (
	ErrLoginDecode          = errors.New("error login parsing user input")
	ErrLoginBusiness        = errors.New("error login from business layer")
	ErrLoginGenAuthToken    = errors.New("error login generating auth token")
	ErrLoginGenRefreshToken = errors.New("error login generating refresh token")

	ErrLogoutDecode               = errors.New("error logout parsing user input")
	ErrLogoutReadRefreshToken     = errors.New("error logout reading refresh token")
	ErrLogoutValidateRefreshToken = errors.New("error logout validating refresh token")
	ErrLogoutBusiness             = errors.New("error logout from business layer")
	ErrLogoutRevokeToken          = errors.New("error logout revoking token")

	ErrChangePassDecode               = errors.New("error change password parsing user input")
	ErrChangePassReadRefreshToken     = errors.New("error change password reading refresh token")
	ErrChangePassValidateRefreshToken = errors.New("error change password validating refresh token")
	ErrChangePassBusiness             = errors.New("error change password from business layer")
	ErrChangePassRevokeToken          = errors.New("error change password revoking token")
)
