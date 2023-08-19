package messages

type ResetEmail struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	ResetToken string `json:"reset_token"`
}
