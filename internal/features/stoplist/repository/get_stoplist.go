package feature_repository_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"context"
	"fmt"
	"time"
)

func (r *RepositoryStopList) GetStopList(
	ctx context.Context,
	limit *int,
	offset *int,
) (core_domain.StopList, error) {
	ctx, close := context.WithTimeout(ctx, 300*time.Millisecond)
	defer close()

	list, err := r.rds.SMembers(ctx, "stoplist:queries").Result()
	if err != nil {
		return core_domain.StopList{}, fmt.Errorf("get stop list: %w", err)
	}

	listDomain := core_domain.NewStopList(list)
	if listDomain == nil {
		return core_domain.StopList{}, nil
	}

	return *listDomain, nil
}
