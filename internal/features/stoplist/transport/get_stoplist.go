package feature_transport_stoplist

import (
	"RWBDwmoTask/internal/core/transport/http/response"
	core_http_utils "RWBDwmoTask/internal/core/transport/http/utils"
	"net/http"
)

func (t *TransportStopList) GetStopList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//log := core_logger.FromContext(ctx)
	responseHandler := response.NewHandlerResponse(nil, w)

	limit, err := core_http_utils.GetIntQueryParm(r, "limit")
	if err != nil {
		responseHandler.ErrorResponse(err, "limit is required")
		return
	}
	offset, err := core_http_utils.GetIntQueryParm(r, "offset")
	if err != nil {
		responseHandler.ErrorResponse(err, "offset is required")
		return
	}

	list, err := t.serviceStopList.GetStopList(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	responseHandler.JSONResponseHandler(http.StatusOK, list)
}
