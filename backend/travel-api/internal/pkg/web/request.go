package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

type validator interface {
	Validate() error
}

func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

func Decode(r *http.Request, val any) error {

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(val); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("validation error: %w", err)
		}
	}

	return nil
}
