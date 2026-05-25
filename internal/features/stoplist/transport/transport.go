package feature_transport_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	core_server "RWBDwmoTask/internal/core/transport/server"
	"context"
	"net/http"
)

type ServiceStopList interface {
	GetStopList(
		ctx context.Context,
		limit *int,
		offset *int,
	) (*core_domain.StopList, error)

	DeleteStopList(
		ctx context.Context,
		id string,
	) error

	AddStopList(
		ctx context.Context,
		item string,
	) error
}

type TransportStopList struct {
	serviceStopList ServiceStopList
}

func NewTransportStopList(
	serviceStopList ServiceStopList,
) *TransportStopList {
	return &TransportStopList{
		serviceStopList: serviceStopList,
	}
}

func (t *TransportStopList) Router() []core_server.Route {
	return []core_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/stoplist",
			Handler: t.GetStopList,
		},
		{
			Method:  http.MethodPost,
			Path:    "/stoplist/{key}",
			Handler: t.AddStopList,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/stoplist/{key}",
			Handler: t.DeleteStopList,
		},
	}
}
