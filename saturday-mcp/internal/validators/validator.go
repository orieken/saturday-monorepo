package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator for request validation
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	v := validator.New()
	
	// Register custom validators
	v.RegisterValidation("validName", validateName)
	v.RegisterValidation("validSelector", validateSelector)
	
	return &Validator{
		validate: v,
	}
}

// Validate validates a struct using validation tags
func (v *Validator) Validate(data interface{}) error {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	// Convert validation errors to user-friendly format
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return NewValidationError(validationErrors)
	}

	return err
}

// ValidateField validates a single field
func (v *Validator) ValidateField(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// validateName ensures names follow naming conventions
func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if len(name) == 0 {
		return false
	}

	// Name should start with a letter and contain only alphanumeric and underscores
	for i, r := range name {
		if i == 0 {
			if !isLetter(r) {
				return false
			}
		} else {
			if !isLetter(r) && !isDigit(r) && r != '_' && r != '-' {
				return false
			}
		}
	}

	return true
}

// validateSelector ensures CSS selectors are valid
func validateSelector(fl validator.FieldLevel) bool {
	selector := fl.Field().String()
	if len(selector) == 0 {
		return false
	}

	// Basic selector validation - must start with # . [ or be a tag name
	firstChar := rune(selector[0])
	return firstChar == '#' || firstChar == '.' || firstChar == '[' || isLetter(firstChar)
}

// isLetter checks if a rune is a letter
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// isDigit checks if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// ValidationError represents validation errors
type ValidationError struct {
	Errors map[string]string
}

// NewValidationError creates a ValidationError from validator errors
func NewValidationError(errs validator.ValidationErrors) *ValidationError {
	errors := make(map[string]string)
	
	for _, err := range errs {
		field := err.Field()
		tag := err.Tag()
		
		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "min":
			message = fmt.Sprintf("%s must have at least %s items", field, err.Param())
		case "url":
			message = fmt.Sprintf("%s must be a valid URL", field)
		case "oneof":
			message = fmt.Sprintf("%s must be one of: %s", field, err.Param())
		case "validName":
			message = fmt.Sprintf("%s must be a valid identifier", field)
		case "validSelector":
			message = fmt.Sprintf("%s must be a valid CSS selector", field)
		default:
			message = fmt.Sprintf("%s failed validation: %s", field, tag)
		}
		
		errors[strings.ToLower(field)] = message
	}
	
	return &ValidationError{Errors: errors}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	var messages []string
	for _, msg := range e.Errors {
		messages = append(messages, msg)
	}
	return strings.Join(messages, "; ")
}

// GetErrors returns the validation errors map
func (e *ValidationError) GetErrors() map[string]string {
	return e.Errors
}
