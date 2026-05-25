package feature_service_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"context"
)

func (s *ServiceStopList) GetStopList(
	ctx context.Context,
	limit *int,
	offset *int,
) (*core_domain.StopList, error) {
	return core_domain.NewStopList(s.stopList.Items()), nil
}
