package feature_service_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"context"
)

type RepositoryStopList interface {
	GetStopList(
		ctx context.Context,
		limit *int,
		offset *int,
	) (core_domain.StopList, error)

	DeleteStopList(
		ctx context.Context,
		id string,
	) error

	AddStopList(
		ctx context.Context,
		item string,
	) error
}

type ServiceStopList struct {
	repository RepositoryStopList
	stopList   *core_domain.StopList
}

func NewServiceStopList(
	repository RepositoryStopList,
	stopList *core_domain.StopList,
) *ServiceStopList {
	return &ServiceStopList{
		repository: repository,
		stopList:   stopList,
	}
}
