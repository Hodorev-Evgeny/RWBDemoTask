package feature_service_stoplist

import (
	core_errors "RWBDwmoTask/internal/core/errors"
	"context"
	"fmt"
)

func (s *ServiceStopList) AddStopList(
	ctx context.Context,
	item string,
) error {
	if item == "" {
		return fmt.Errorf("item is empty: %w", core_errors.ErrorValidation)
	}

	ans := s.repository.AddStopList(ctx, item)
	if ans != nil {
		return fmt.Errorf("error add stoplist: %w", ans)
	}

	s.stopList.Add(item)
	return nil
}
