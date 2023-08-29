package mail

type Letter struct {
	To      string
	Name    string
	Subject string
	Header  string
	Body    string
	Token   string
}

type Message struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	ResetToken string `json:"reset_token"`
}
