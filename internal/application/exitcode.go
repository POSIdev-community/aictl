package application

import (
	"errors"

	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
)

const (
	ExitCodeSuccess    = 0
	ExitCodeValidation = 1
	ExitCodeAPI        = 2
	ExitCodeUnknown    = -1
)

func mapExitCode(err error) (exitCode int, errorMessage string) {
	var (
		validationErr         *validation.Error
		validationFieldErr    *validation.FieldError
		validationRequiredErr *validation.RequiredError
		validationInvalidErr  *validation.InvalidError

		failErr            *apperror.FailError
		emptyResponseError *apperror.EmptyResponseError
		authenticationErr  *apperror.AuthenticationError
		authorizationErr   *apperror.AuthorizationError
		badRequestErr      *apperror.BadRequestError
		notFoundErr        *apperror.NotFoundError
		notFoundByIdErr    *apperror.NotFoundByIdError
		unknownErr         *apperror.UnknownResponseError
		serverResponseErr  *apperror.ServerResponseError

		badRequestApiErrorModelErr *apperror.BadRequestApiErrorModelError
		unknownApiErrorModelErr    *apperror.UnknownApiErrorModelError
	)

	switch {
	case errors.As(err, &validationErr):
		return ExitCodeValidation, validationErr.Error()
	case errors.As(err, &validationFieldErr):
		return ExitCodeValidation, validationFieldErr.Error()
	case errors.As(err, &validationRequiredErr):
		return ExitCodeValidation, validationRequiredErr.Error()
	case errors.As(err, &validationInvalidErr):
		return ExitCodeValidation, validationInvalidErr.Error()

	case errors.As(err, &failErr):
		return ExitCodeValidation, failErr.Error()
	case errors.As(err, &emptyResponseError):
		return ExitCodeValidation, emptyResponseError.Error()
	case errors.As(err, &authenticationErr):
		return ExitCodeAPI, authenticationErr.Error()
	case errors.As(err, &authorizationErr):
		return ExitCodeAPI, authorizationErr.Error()
	case errors.As(err, &badRequestErr):
		return ExitCodeAPI, badRequestErr.Error()
	case errors.As(err, &notFoundErr):
		return ExitCodeAPI, notFoundErr.Error()
	case errors.As(err, &notFoundByIdErr):
		return ExitCodeAPI, notFoundByIdErr.Error()
	case errors.As(err, &unknownErr):
		return ExitCodeAPI, unknownErr.Error()
	case errors.As(err, &serverResponseErr):
		return ExitCodeAPI, serverResponseErr.Error()
	case errors.As(err, &badRequestApiErrorModelErr):
		return ExitCodeAPI, badRequestApiErrorModelErr.Error()
	case errors.As(err, &unknownApiErrorModelErr):
		return ExitCodeAPI, unknownApiErrorModelErr.Error()

	default:
		return ExitCodeUnknown, err.Error()
	}
}
