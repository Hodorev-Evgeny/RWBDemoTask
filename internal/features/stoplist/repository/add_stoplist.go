package feature_repository_stoplist

import "context"

func (r *RepositoryStopList) AddStopList(
	ctx context.Context,
	item string,
) error {
	return r.rds.SAdd(ctx, "stop_list", item).Err()
}
