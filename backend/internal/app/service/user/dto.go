package user

import "time"

type UserResponseDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateCreated time.Time `json:"date_created"`
}

type NewUserDTO struct {
	// TODO: validation
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

type UpdateUserDTO struct {
	// TODO: validation
	Name            *string `json:"name"`
	Email           *string `json:"email" validate:"omitempty,email"`
	Password        *string `json:"password" validate:"omitempty"`
	PasswordConfirm *string `json:"password_confirm" validate:"eqfield=Password"`
}
