package errs

import (
	"fmt"
	"strings"
)

type BadRequestApiErrorModelError struct {
	errorCode string
	details   map[string]*string
}

func (e *BadRequestApiErrorModelError) Error() string {
	result := make([]string, 0, len(e.details))
	for key, value := range e.details {
		str := fmt.Sprintf("%s: %s", key, *value)
		result = append(result, str)
	}

	if len(result) == 0 {
		return fmt.Sprintf("Bad request: '%s'", e.errorCode)
	} else {
		return fmt.Sprintf("Bad request: '%s'\n%s", e.errorCode, strings.Join(result, "\n"))
	}
}

func NewBadRequestApiErrorModelError(errorCode string, details map[string]*string) *BadRequestApiErrorModelError {
	return &BadRequestApiErrorModelError{errorCode, details}
}

type UnknownApiErrorModelError struct {
	statusCode int
	errorCode  string
	details    map[string]*string
}

func (e *UnknownApiErrorModelError) Error() string {
	result := make([]string, 0, len(e.details))
	for key, value := range e.details {
		str := fmt.Sprintf("%s: %s", key, *value)
		result = append(result, str)
	}

	if len(result) == 0 {
		return fmt.Sprintf("Unknown error. Status code: %d; %s", e.statusCode, e.errorCode)
	} else {
		return fmt.Sprintf("Unknown error. Status code: %d; %s\n%s", e.statusCode, e.errorCode, strings.Join(result, "\n"))
	}
}

func NewUnknownApiErrorModelError(statusCode int, errorCode string, details map[string]*string) *UnknownApiErrorModelError {
	return &UnknownApiErrorModelError{
		statusCode,
		errorCode,
		details,
	}
}

func CheckApiErrorModel(statusCode int, errorCode string, details map[string]*string) error {
	switch {
	case statusCode < 400:
		return nil
	case statusCode == 400:
		return NewBadRequestApiErrorModelError(errorCode, details)
	case statusCode == 401:
		return NewAuthenticationError()
	case statusCode == 403:
		return NewAuthorizationError()
	case statusCode == 404:
		return NewNotFoundError(errorCode)
	default:
		return NewUnknownApiErrorModelError(statusCode, errorCode, details)
	}
}
