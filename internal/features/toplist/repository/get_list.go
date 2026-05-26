package feature_repository_toplist

import (
	core_domain "RWBDwmoTask/internal/core/storage"
	"context"
	"time"
)

func (r *NatsRepository) GetList(
	ctx context.Context,
	limit int,
) ([]core_domain.TopItem, error) {
	ctx, close := context.WithTimeout(ctx, 5*time.Second)
	defer close()

	list := r.storege.GetTop(limit)
	return list, nil
}
