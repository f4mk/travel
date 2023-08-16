package web

import (
	"encoding/json"
	"errors"
	"net/http"
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
