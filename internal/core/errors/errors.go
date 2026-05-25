package core_errors

import (
	"errors"
)

var (
	ErrorNotFoud      = errors.New("Not foud")
	ErrorBadRequest   = errors.New("Bad request")
	ErrorUnauthorized = errors.New("Unauthorized")
	ErrorValidation   = errors.New("Validation error")
)
