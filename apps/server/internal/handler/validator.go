package handler

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"omni-pixel/internal/apperr"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Report the JSON field name in validation errors so clients see
	// `email` instead of `Email`.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})
	return v
}

// validateStruct runs the validator and converts any validator.ValidationErrors
// into an apperr.ErrValidation carrying a details map keyed by JSON field name.
func validateStruct(v any) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}
	var verr validator.ValidationErrors
	if !asValidationErrors(err, &verr) {
		return apperr.ErrValidation.WithCause(err)
	}
	details := make(map[string]any, len(verr))
	for _, fe := range verr {
		details[fe.Field()] = fe.Tag()
	}
	return apperr.ErrValidation.WithDetails(details).WithCause(err)
}

func asValidationErrors(err error, out *validator.ValidationErrors) bool {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}
	*out = ve
	return true
}
