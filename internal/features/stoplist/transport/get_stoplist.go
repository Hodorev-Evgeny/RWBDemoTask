package feature_transport_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"RWBDwmoTask/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type StopListResponse struct {
	List []string `json:"list"`
}

func (t *TransportStopList) GetStopList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseHandler := response.NewHandlerResponse(nil, w)

	list, err := t.serviceStopList.GetStopList(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, err.Error())
		return
	}

	req := DomainFromResponse(list)
	fmt.Println(req)
	responseHandler.JSONResponseHandler(http.StatusOK, req)
}

func DomainFromResponse(
	list *core_domain.StopList,
) StopListResponse {
	return StopListResponse{
		List: list.Items(),
	}
}
