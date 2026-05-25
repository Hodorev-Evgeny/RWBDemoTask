package core_http_utils

import (
	core_errors "RWBDwmoTask/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetValuePathInt(r *http.Request, key string) (int, error) {
	value := r.PathValue(key)
	if value == "" {
		return 0, fmt.Errorf(`path "%s" is required`, key)
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid value for path %s: %e", key, core_errors.ErrorValidation)
	}

	return valueInt, nil
}

func GetValuePathString(r *http.Request, key string) (string, error) {
	value := r.PathValue(key)
	if value == "" {
		return "", fmt.Errorf(`path "%s" is required`, key)
	}

	return value, nil
}
