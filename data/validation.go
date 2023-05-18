package data

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// ValidationError wraps the validators FieldError so we do not
// expose this out
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// Validation contains
type Validation struct {
	validate *validator.Validate
}

// NewValidation creates a new Validation type
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	//errs := v.validate.Struct(i).(validator.ValidationErrors)
	//if errs == nil {
	//	return nil
	//}
	//
	//if len(errs) == 0 {
	//	return nil
	//}
	errs := v.validate.Struct(i)
	if errs == nil {
		return nil
	}

	validationErrors, ok := errs.(validator.ValidationErrors)
	if !ok {
		// Handle the case when the returned value is not of type validator.ValidationErrors
		// You can choose to return an error or handle it according to your application's logic
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range validationErrors {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-absd-sdsd
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true
}
