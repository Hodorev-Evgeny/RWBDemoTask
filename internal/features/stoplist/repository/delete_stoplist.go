package feature_repository_stoplist

import (
	"context"
)

func (r *RepositoryStopList) DeleteStopList(
	ctx context.Context,
	id string,
) error {
	return r.rds.SRem(ctx, "stoplist:queries", id).Err()
}
