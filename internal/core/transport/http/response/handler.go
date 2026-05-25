package response

import (
	core_errors "RWBDwmoTask/internal/core/errors"
	core_logger "RWBDwmoTask/internal/core/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HandlerResponse struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHandlerResponse(logger *core_logger.Logger, rw http.ResponseWriter) *HandlerResponse {
	return &HandlerResponse{
		log: logger,
		rw:  rw,
	}
}

func (h *HandlerResponse) JSONResponseHandler(statuscode int, data any) {
	h.rw.WriteHeader(statuscode)

	if err := json.NewEncoder(h.rw).Encode(data); err != nil {
		fmt.Println("Failed to encode response", err)
	}
}

func (h *HandlerResponse) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logfunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrorBadRequest):
		statusCode = http.StatusBadRequest

	case errors.Is(err, core_errors.ErrorUnauthorized):
		statusCode = http.StatusUnauthorized

	case errors.Is(err, core_errors.ErrorValidation):
		statusCode = http.StatusUnprocessableEntity

	case errors.Is(err, core_errors.ErrorNotFoud):
		statusCode = http.StatusNotFound

	default:
		statusCode = http.StatusInternalServerError
	}

	logfunc(msg, zap.Error(err))

	h.errorResponse(err, msg, statusCode)
}

func (h *HandlerResponse) PanicResponse(p any, msg string) {
	statuscode := http.StatusInternalServerError
	err := fmt.Errorf("unexepted punic:", p)

	h.errorResponse(err, msg, statuscode)
}

func (h *HandlerResponse) errorResponse(err error, msg string, status int) {
	h.rw.WriteHeader(status)

	response := ErrorResponse{
		Error:   msg,
		Massage: msg,
	}

	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		fmt.Println("Failed to encode response", err)
	}
}
