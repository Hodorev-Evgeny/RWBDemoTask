package feature_service_toplist

import (
	core_domain "RWBDwmoTask/internal/core/storage"
	"context"
	"fmt"
)

func (s *NatsService) GetList(
	ctx context.Context,
	limit int,
) ([]core_domain.TopItem, error) {
	// add stop list

	list, err := s.repository.GetList(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("GetList on service: %w", err)
	}
	return list, nil
}
