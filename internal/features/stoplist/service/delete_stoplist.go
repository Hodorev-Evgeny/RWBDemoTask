package feature_service_stoplist

import (
	core_errors "RWBDwmoTask/internal/core/errors"
	"context"
	"fmt"
)

func (s *ServiceStopList) DeleteStopList(
	ctx context.Context,
	id string,
) error {
	if id == "" {
		return fmt.Errorf("id is required:%w", core_errors.ErrorValidation)
	}

	ans := s.repository.DeleteStopList(ctx, id)
	if ans != nil {
		return fmt.Errorf("error deleate stoplist: %w", ans)
	}

	s.stopList.Remove(id)
	return nil
}
