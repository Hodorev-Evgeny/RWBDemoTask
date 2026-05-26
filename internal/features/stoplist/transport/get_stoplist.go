package feature_transport_stoplist

import (
	"RWBDwmoTask/internal/core/transport/http/response"
	"net/http"
)

func (t *TransportStopList) GetStopList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseHandler := response.NewHandlerResponse(nil, w)

	list, err := t.serviceStopList.GetStopList(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponseHandler(http.StatusOK, list)
}
