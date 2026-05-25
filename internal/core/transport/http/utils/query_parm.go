package core_http_utils

import (
	core_errors "RWBDwmoTask/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetIntQueryParm(r *http.Request, key string) (*int, error) {
	value := r.URL.Query().Get(key)

	if value == "" {
		return nil, nil
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value for value %s: %e", key, core_errors.ErrorValidation)
	}

	return &valueInt, nil
}
