package auth

type User struct {
	ID           string `db:"user_id"`
	Email        string `db:"email"`
	PasswordHash []byte `db:"password_hash"`
}

type Login struct {
	// TODO: validation
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// TODO: validation
type Logout struct {
	Email string `json:"email" validate:"required"`
}

// TODO: validation
type Refresh struct {
	Token string `json:"token" validate:"required"`
}

// TODO: validation
type PasswordReset struct {
	Email string `json:"email" validate:"required"`
}

type StoreToken struct {
	Token     string `db:"token"`
	UserID    string `db:"user_id"`
	ExpiresAt int64  `db:"expires_at"`
	CreatedAt int64  `db:"created_at"`
}
