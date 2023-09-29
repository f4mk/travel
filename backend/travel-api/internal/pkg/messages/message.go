package messages

type MessageType int

const (
	ResetPassword MessageType = iota
	RegisterVerify
)

type Message struct {
	ID    string      `json:"id"`
	Email string      `json:"email"`
	Name  string      `json:"name"`
	Token string      `json:"token"`
	Type  MessageType `json:"type"`
}
