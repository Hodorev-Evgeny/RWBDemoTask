package feature_repository_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"context"
	"fmt"
	"time"
)

func (r *RepositoryStopList) GetStopList(
	ctx context.Context,
) (*core_domain.StopList, error) {
	ctx, close := context.WithTimeout(ctx, 300*time.Millisecond)
	defer close()

	list, err := r.rds.GetStoplist(ctx)
	if err != nil {
		return nil, fmt.Errorf("get stop list: %w", err)
	}

	listDomain := core_domain.NewStopList(list)
	if listDomain == nil {
		return nil, nil
	}

	return listDomain, nil
}
