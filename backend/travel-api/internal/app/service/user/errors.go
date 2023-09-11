package user

import "errors"

var (
	ErrGetUsersBusiness = errors.New("error get users from business layer")
	ErrValidateUserID   = errors.New("error get user validate uuid")
	ErrGetUserBusiness  = errors.New("error get user from business layer")

	ErrCreateValidate    = errors.New("error create parsing user input")
	ErrCreateBusiness    = errors.New("error create user from business layer")
	ErrCreateSendMessage = errors.New("error create user send message")

	ErrUpdateValidateUUID = errors.New("error update user validate uuid")
	ErrUpdateValidate     = errors.New("error update parsing user input")
	ErrUpdateBusiness     = errors.New("error update user from business layer")

	ErrVerifyValidate = errors.New("error update parsing user input")
	ErrVerifyBusiness = errors.New("error update user from business layer")

	ErrDeleteValidateUUID      = errors.New("error delete user validate uuid")
	ErrDeleteBusiness          = errors.New("error delete user from business layer")
	ErrDeleteStoreTokenVersion = errors.New("error delete user storing token version")
)
