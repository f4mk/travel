package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {

	// TODO: think what to do with this
	_ = SetStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {

	var re ResponseError
	var status int

	switch {
	case IsRequestError(err):
		reqErr := GetRequestError(err)
		if IsFieldErrors(reqErr.Err) {
			fieldErrors := GetFieldErrors(reqErr.Err)

			re = ResponseError{
				Error:  "data validation error",
				Fields: fieldErrors.Fields(),
			}
			status = reqErr.Status

		} else {
			re = ResponseError{
				Error: reqErr.Error(),
			}
			status = reqErr.Status
		}

	default:
		re = ResponseError{
			Error: http.StatusText(http.StatusInternalServerError),
		}
		status = http.StatusInternalServerError
	}

	if err := Respond(ctx, w, re, status); err != nil {
		return err
	}

	return nil
}
