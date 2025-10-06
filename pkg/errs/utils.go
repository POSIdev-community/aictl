package errs

import "errors"

var (
	validationErr      *ValidationError
	validationFieldErr *ValidationFieldError

	nilResponseError  *NilResponseError
	authenticationErr *AuthenticationError
	authorizationErr  *AuthorizationError
	badRequestErr     *BadRequestError
	notFoundErr       *NotFoundError
	notFoundByIdErr   *NotFoundByIdError
	unknownErr        *UnknownResponseError
	serverResponseErr *ServerResponseError

	badRequestApiErrorModelErr *BadRequestApiErrorModelError
	unknownApiErrorModelErr    *UnknownApiErrorModelError
)

func MapExitCode(err error) (exitCode int, errorMessage string) {
	switch {
	case errors.As(err, &validationErr):
		exitCode = 1
		errorMessage = validationErr.Error()
	case errors.As(err, &validationFieldErr):
		exitCode = 1
		errorMessage = validationFieldErr.Error()

	case errors.As(err, &nilResponseError):
		exitCode = 1
		errorMessage = nilResponseError.Error()
	case errors.As(err, &authenticationErr):
		exitCode = 2
		errorMessage = authenticationErr.Error()
	case errors.As(err, &authorizationErr):
		exitCode = 2
		errorMessage = authorizationErr.Error()
	case errors.As(err, &badRequestErr):
		exitCode = 2
		errorMessage = badRequestErr.Error()
	case errors.As(err, &notFoundErr):
		exitCode = 2
		errorMessage = notFoundErr.Error()
	case errors.As(err, &notFoundByIdErr):
		exitCode = 2
		errorMessage = notFoundByIdErr.Error()
	case errors.As(err, &unknownErr):
		exitCode = 2
		errorMessage = unknownErr.Error()
	case errors.As(err, &serverResponseErr):
		exitCode = 2
		errorMessage = serverResponseErr.Error()

	case errors.As(err, &badRequestApiErrorModelErr):
		exitCode = 2
		errorMessage = badRequestApiErrorModelErr.Error()
	case errors.As(err, &unknownApiErrorModelErr):
		exitCode = 2
		errorMessage = unknownApiErrorModelErr.Error()

	default:
		exitCode = -1
		errorMessage = err.Error()
	}

	return exitCode, errorMessage
}
