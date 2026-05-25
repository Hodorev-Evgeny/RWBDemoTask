package feature_transport_toplist

import (
	"encoding/json"
	"net/http"
)

func (t *TransportTopList) GetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//log := core_logger.FromContext(ctx)

	list, err := t.service.GetList(ctx, 5)
	if err != nil {
		//log.Error("GetList error", zap.Error(err))
		http.ResponseWriter(w).WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ResponseWriter(w).WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		//log.Error("GetList error", zap.Error(err))
		return
	}
}
