package core_http_utils

import (
	core_errors "RWBDwmoTask/internal/core/errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var v = validator.New()

type customValidator interface {
	Validate() error
}

func DecodeJSON(data any, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	var err error

	value, ok := data.(customValidator)
	if ok {
		err = value.Validate()
	} else {
		err = v.Struct(data)
	}

	if err != nil {
		return fmt.Errorf("invalid argument: %w, %w",
			err,
			core_errors.ErrorBadRequest,
		)
	}

	return nil
}
