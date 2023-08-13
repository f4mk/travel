package auth

type LoginUser struct {
	// TODO: validation
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
