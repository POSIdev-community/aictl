package errs

import "fmt"

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s", e.Message)
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

type ValidationFieldError struct {
	Field   string
	Message string
}

func (e *ValidationFieldError) Error() string {
	return fmt.Sprintf("Validation error on field %s: %s", e.Field, e.Message)
}

func NewValidationFieldError(field, message string) *ValidationFieldError {
	return &ValidationFieldError{
		Field:   field,
		Message: message,
	}
}

type ValidationRequiredError struct {
	Field string
}

func (e *ValidationRequiredError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': param is required", e.Field)
}

func NewValidationRequiredError(field string) *ValidationRequiredError {
	return &ValidationRequiredError{
		Field: field,
	}
}

type ValidationInvalidError struct {
	Field string
}

func (e *ValidationInvalidError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': param is invalid", e.Field)
}

func NewValidationInvalidError(field string) *ValidationInvalidError {
	return &ValidationInvalidError{
		Field: field,
	}
}
