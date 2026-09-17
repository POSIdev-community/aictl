package validation

import "fmt"

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Validation error: %s", e.Message)
}

func NewError(message string) *Error {
	return &Error{Message: message}
}

type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("Validation error on field %s: %s", e.Field, e.Message)
}

func NewFieldError(field, message string) *FieldError {
	return &FieldError{Field: field, Message: message}
}

type RequiredError struct {
	Field string
}

func (e *RequiredError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': param is required", e.Field)
}

func NewRequiredError(field string) *RequiredError {
	return &RequiredError{Field: field}
}

type InvalidError struct {
	Field string
}

func (e *InvalidError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': param is invalid", e.Field)
}

func NewInvalidError(field string) *InvalidError {
	return &InvalidError{Field: field}
}

// MessageError is a validation failure whose Error() is the message as-is (no "Validation error:" prefix).
type MessageError struct {
	Message string
}

func (e *MessageError) Error() string {
	return e.Message
}

func NewMessageError(message string) *MessageError {
	return &MessageError{Message: message}
}
