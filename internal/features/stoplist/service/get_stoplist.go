package feature_service_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"context"
	"fmt"
)

func (s *ServiceStopList) GetStopList(
	ctx context.Context,
) (*core_domain.StopList, error) {
	list, err := s.repository.GetStopList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get stop list on service: %w", err)
	}

	return list, nil
}
