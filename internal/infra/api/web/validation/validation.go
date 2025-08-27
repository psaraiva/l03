package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"l03/configuration/rest_err"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validator_en "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate   = validator.New()
	translator ut.Translator
)

func init() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		enTranslate := ut.New(en, en)
		translator, _ = enTranslate.GetTranslator("en")
		validator_en.RegisterDefaultTranslations(value, translator)
	}
}

func ValidateErr(validations_err error) *rest_err.RestError {
	var jsonErr *json.UnmarshalTypeError
	var jsonValidation validator.ValidationErrors

	if errors.As(validations_err, &jsonErr) {
		return rest_err.NewBadRequestError(fmt.Sprintf("Invalid field type, '%s'", jsonErr.Field))
	} else if errors.As(validations_err, &jsonValidation) {
		errorCauses := []rest_err.Causes{}

		for _, e := range validations_err.(validator.ValidationErrors) {
			errorCauses = append(errorCauses, rest_err.Causes{
				Field:   e.Field(),
				Message: e.Translate(translator),
			})
		}
		return rest_err.NewBadRequestError("Some fields are invalid", errorCauses...)
	}
	return rest_err.NewBadRequestError("Error trying to convert fields")
}
