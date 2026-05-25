package feature_transport_stoplist

import (
	"RWBDwmoTask/internal/core/transport/http/response"
	core_http_utils "RWBDwmoTask/internal/core/transport/http/utils"
	"net/http"
)

func (t *TransportStopList) AddStopList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseHandler := response.NewHandlerResponse(nil, w)

	key, err := core_http_utils.GetValuePathString(r, "key")
	if err != nil {
		responseHandler.ErrorResponse(err, "error getting key")
		return
	}

	err = t.serviceStopList.AddStopList(ctx, key)
	if err != nil {
		responseHandler.ErrorResponse(err, "error adding stop list")
		return
	}

	responseHandler.JSONResponseHandler(http.StatusCreated, key)
}
