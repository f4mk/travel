package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
)

func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	ctx, span := AddSpan(ctx, "web.response", attribute.Int("status", statusCode))
	defer span.End()
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

func RespondRaw(
	ctx context.Context,
	w http.ResponseWriter,
	data io.Reader,
	statusCode int,
	ctype string,
) error {
	ctx, span := AddSpan(ctx, "web.response-raw", attribute.Int("status", statusCode))
	defer span.End()
	_ = SetStatusCode(ctx, statusCode)

	w.Header().Set("Content-Type", ctype)
	w.WriteHeader(statusCode)
	_, err := io.Copy(w, data)
	if err != nil {
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
	case IsTimeoutError(err):
		re = ResponseError{
			Error: "request timeout",
		}
		status = http.StatusRequestTimeout
	default:
		re = ResponseError{
			Error: http.StatusText(http.StatusInternalServerError),
		}
		status = http.StatusInternalServerError
	}

	return Respond(ctx, w, re, status)
}
