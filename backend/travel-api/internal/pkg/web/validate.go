package web

import (
	"reflect"
	"strings"

	"github.com/google/uuid"

	errValidator "github.com/go-playground/validator/v10"
)

var validate *errValidator.Validate

func init() {

	validate = errValidator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

}

// Check validates the provided model against it's declared tags.
func Check(val any) error {
	if err := validate.Struct(val); err != nil {

		verrors, ok := err.(errValidator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Error(),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}

func ValidateUUID(id string) error {
	_, err := uuid.Parse(id)
	if err != nil {

		return err
	}
	return nil
}
