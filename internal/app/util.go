package app

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(r interface{}) error {
	if err := validate.Struct(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid request")
		}

		for _, err := range err.(validator.ValidationErrors) {
			var errMsg string
			switch err.Tag() {
			case "required":
				errMsg = fmt.Sprintf("%s field is required", err.Field())
			case "email":
				errMsg = fmt.Sprintf("%s field is not valid", err.Field())
			case "min":
				errMsg = fmt.Sprintf("%s field must be at least %s", err.Field(), err.Param())
			case "max":
				errMsg = fmt.Sprintf("%s field must be at most %s", err.Field(), err.Param())
			case "gte":
				errMsg = fmt.Sprintf("%s field must be greater than or equal to %s", err.Field(), err.Param())
			case "lte":
				errMsg = fmt.Sprintf("%s field must be less than or equal to %s", err.Field(), err.Param())
			case "eqfield":
				errMsg = fmt.Sprintf("%s field must be equal to %s", err.Field(), err.Param())
			case "gtfield":
				errMsg = fmt.Sprintf("%s field must be greater than %s", err.Field(), err.Param())
			case "gtefield":
				errMsg = fmt.Sprintf("%s field must be greater than or equal to %s", err.Field(), err.Param())
			case "ltfield":
				errMsg = fmt.Sprintf("%s field must be less than %s", err.Field(), err.Param())
			case "ltefield":
				errMsg = fmt.Sprintf("%s field must be less than or equal to %s", err.Field(), err.Param())
			case "nefield":
				errMsg = fmt.Sprintf("%s field must not be equal to %s", err.Field(), err.Param())
			case "uniquefield":
				errMsg = fmt.Sprintf("%s field must be unique", err.Field())
			case "numeric":
				errMsg = fmt.Sprintf("%s field must be numeric", err.Field())
			case "alphanum":
				errMsg = fmt.Sprintf("%s field must be alphanumeric", err.Field())
			default:
				errMsg = fmt.Sprintf("%s field is invalid (%s)", err.Field(), err.Tag())
			}

			if errMsg != "" {
				return NewBadRequestError(errors.New(errMsg))
			}
		}

		return err
	}

	return nil
}
