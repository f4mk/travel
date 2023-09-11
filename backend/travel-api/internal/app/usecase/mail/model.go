package mail

type Letter struct {
	To      string
	Name    string
	Subject string
	Header  string
	Body    string
	Token   string
}

type MessageReset struct {
	Email      string
	Name       string
	ResetToken string
}

type MessageVerify struct {
	Email       string
	Name        string
	VerifyToken string
}
