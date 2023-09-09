// Package auth provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package auth

import (
	"time"
)

// ChangePassword defines model for ChangePassword.
type ChangePassword struct {
	// Password user new password
	Password string `json:"password" validate:"required,gte=8"`

	// PasswordConfirm user new password confirm
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`

	// PasswordOld user new password
	PasswordOld string `json:"password_old" validate:"required"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Error error message
	Error  string             `json:"error"`
	Fields *map[string]string `json:"fields,omitempty"`
}

// LoginUser defines model for LoginUser.
type LoginUser struct {
	// Email user email
	Email string `json:"email" validate:"required,email"`

	// Password user password
	Password string `json:"password" validate:"required"`
}

// NewPassword defines model for NewPassword.
type NewPassword struct {
	// Password user new password
	Password string `json:"password" validate:"required,gte=8"`

	// PasswordConfirm user new password confirm
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

// OldPassword defines model for OldPassword.
type OldPassword struct {
	// PasswordOld user new password
	PasswordOld string `json:"password_old" validate:"required"`
}

// ResetPassword defines model for ResetPassword.
type ResetPassword struct {
	// Email user email
	Email string `json:"email" validate:"required,email"`
}

// ResetToken defines model for ResetToken.
type ResetToken struct {
	// Token password reset secret token
	Token string `json:"token" validate:"required"`
}

// SubmitResetPassword defines model for SubmitResetPassword.
type SubmitResetPassword struct {
	// Password user new password
	Password string `json:"password" validate:"required,gte=8"`

	// PasswordConfirm user new password confirm
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`

	// Token password reset secret token
	Token string `json:"token" validate:"required"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	// DateCreated date created
	DateCreated time.Time `json:"date_created"`

	// Email user's email
	Email string `json:"email"`

	// Id user unique id
	ID string `json:"id"`

	// Name user's name
	Name string `json:"name"`
}

// PostAuthLogoutJSONBody defines parameters for PostAuthLogout.
type PostAuthLogoutJSONBody = map[string]interface{}

// PostAuthLogoutAllJSONBody defines parameters for PostAuthLogoutAll.
type PostAuthLogoutAllJSONBody = map[string]interface{}

// PostAuthRefreshJSONBody defines parameters for PostAuthRefresh.
type PostAuthRefreshJSONBody = map[string]interface{}

// PostAuthValidateJSONBody defines parameters for PostAuthValidate.
type PostAuthValidateJSONBody = map[string]interface{}

// PostAuthLoginJSONRequestBody defines body for PostAuthLogin for application/json ContentType.
type PostAuthLoginJSONRequestBody = LoginUser

// PostAuthLogoutJSONRequestBody defines body for PostAuthLogout for application/json ContentType.
type PostAuthLogoutJSONRequestBody = PostAuthLogoutJSONBody

// PostAuthLogoutAllJSONRequestBody defines body for PostAuthLogoutAll for application/json ContentType.
type PostAuthLogoutAllJSONRequestBody = PostAuthLogoutAllJSONBody

// PostAuthPasswordChangeJSONRequestBody defines body for PostAuthPasswordChange for application/json ContentType.
type PostAuthPasswordChangeJSONRequestBody = ChangePassword

// PostAuthPasswordResetJSONRequestBody defines body for PostAuthPasswordReset for application/json ContentType.
type PostAuthPasswordResetJSONRequestBody = ResetPassword

// PostAuthPasswordResetSubmitJSONRequestBody defines body for PostAuthPasswordResetSubmit for application/json ContentType.
type PostAuthPasswordResetSubmitJSONRequestBody = SubmitResetPassword

// PostAuthRefreshJSONRequestBody defines body for PostAuthRefresh for application/json ContentType.
type PostAuthRefreshJSONRequestBody = PostAuthRefreshJSONBody

// PostAuthValidateJSONRequestBody defines body for PostAuthValidate for application/json ContentType.
type PostAuthValidateJSONRequestBody = PostAuthValidateJSONBody
