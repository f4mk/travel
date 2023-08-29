package mail

type ServeError struct {
	Error   error
	Payload []byte
}
