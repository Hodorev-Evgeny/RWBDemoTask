package feature_transport_toplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	core_server "RWBDwmoTask/internal/core/transport/server"
	"context"
	"net/http"
)

type ServiceTopList interface {
	GetList(
		ctx context.Context,
		limit int,
	) ([]core_domain.TopItem, error)
}

type TransportTopList struct {
	service ServiceTopList
}

func NewTransportTopList(
	service ServiceTopList,
) *TransportTopList {
	return &TransportTopList{
		service: service,
	}
}

func (t *TransportTopList) Router() []core_server.Route {
	return []core_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/toplist",
			Handler: t.GetList,
		},
	}
}
