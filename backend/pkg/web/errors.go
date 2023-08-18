package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"syscall"
)

var (
	ErrNotFound      = errors.New("error not found")
	ErrForbidden     = errors.New("error not allowed")
	ErrAuthFailed    = errors.New("error authentication failed")
	ErrAlreadyExists = errors.New("error already exists")
	ErrGenHash       = errors.New("error generate hash")
	ErrGetClaims     = errors.New("error get claims")
	ErrQueryDB       = errors.New("error query db")
	ErrCritical      = errors.New("error data integrity")
)

type ResponseError struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type RequestError struct {
	Err    error
	Status int
}

type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

type FieldErrors []FieldError

// NewFieldsError creates an fields error.
func NewFieldsError(field string, err error) error {
	return FieldErrors{
		{
			Field: field,
			Err:   err.Error(),
		},
	}
}

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Fields returns the fields that failed validation
func (fe FieldErrors) Fields() map[string]string {
	m := make(map[string]string)
	for _, fld := range fe {
		m[fld.Field] = fld.Err
	}
	return m
}

// IsFieldErrors checks if an error of type FieldErrors exists.
func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

// GetFieldErrors returns a copy of the FieldErrors pointer.
func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}

func (err *Error) Error() string {
	return err.Err.Error()
}

type shutdown struct {
	Message string
}

func NewShutdownError(message string) error {
	return &shutdown{message}
}

func (s shutdown) Error() string {
	return s.Message
}

func IsShutdown(err error) bool {
	if _, ok := err.(*shutdown); ok {
		return true
	}
	return false
}

func IsTimeoutError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}

func GetResponseErrorFromBusiness(err error) error {

	switch {
	case errors.Is(err, ErrNotFound):
		return NewRequestError(
			err,
			http.StatusNotFound,
		)

	case errors.Is(err, ErrForbidden):
		return NewRequestError(
			err,
			http.StatusForbidden,
		)

	case errors.Is(err, ErrAlreadyExists):
		return NewRequestError(
			err,
			http.StatusConflict,
		)
	case errors.Is(err, ErrAuthFailed):
		return NewRequestError(
			err,
			http.StatusUnauthorized,
		)
	default:
		return err
	}
}

func validateShutdown(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}
