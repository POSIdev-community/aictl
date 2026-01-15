package errs

import (
	"fmt"
)

type NilResponseError struct{}

func (e *NilResponseError) Error() string {
	return "Response is nil"
}
func NewNilResponseError() *NilResponseError {
	return &NilResponseError{}
}

type AuthenticationError struct{}

func (e *AuthenticationError) Error() string {
	return "Authentication error"
}
func NewAuthenticationError() *AuthenticationError {
	return &AuthenticationError{}
}

type AuthorizationError struct{}

func (e *AuthorizationError) Error() string {
	return "Authorization error"
}
func NewAuthorizationError() *AuthorizationError {
	return &AuthorizationError{}
}

type BadRequestError struct {
	body string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad Request error: %s", e.body)
}

func NewBadRequestError(body string) *BadRequestError {
	return &BadRequestError{body: body}
}

type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
	}
}

type NotFoundByIdError struct {
	Resource string
	ID       string
}

func (e *NotFoundByIdError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

func NewNotFoundByIdError(resource, id string) *NotFoundByIdError {
	return &NotFoundByIdError{
		Resource: resource,
		ID:       id,
	}
}

type UnknownResponseError struct {
	statusCode int
	body       string
}

func (e *UnknownResponseError) Error() string {
	if e.body == "" {
		return fmt.Sprintf("Unknown response error. Status code: '%d'", e.statusCode)
	} else {
		return fmt.Sprintf("Unknown response error. Status code: '%d'\n%s", e.statusCode, e.body)
	}
}

func NewUnknownResponseError(statusCode int, body string) *UnknownResponseError {
	return &UnknownResponseError{
		statusCode,
		body,
	}
}

type ServerResponseError struct {
	statusCode int
	body       string
}

func (e *ServerResponseError) Error() string {
	if e.body == "" {
		return fmt.Sprintf("Server response error. Status code: '%d'", e.statusCode)
	} else {
		return fmt.Sprintf("Server response error. Status code: '%d'\n%s", e.statusCode, e.body)
	}
}

func NewServerError(statusCode int, body string) *ServerResponseError {
	return &ServerResponseError{
		statusCode,
		body,
	}
}

func CheckResponseErrors(statusCode int, body string, resourceName string) error {
	switch {
	case statusCode < 400:
		return nil
	case statusCode == 400:
		return NewBadRequestError(body)
	case statusCode == 401:
		return NewAuthenticationError()
	case statusCode == 403:
		return NewAuthorizationError()
	case statusCode == 404:
		return NewNotFoundError(resourceName)
	case statusCode < 500:
		return NewUnknownResponseError(statusCode, body)
	default:
		return NewServerError(statusCode, body)
	}
}
