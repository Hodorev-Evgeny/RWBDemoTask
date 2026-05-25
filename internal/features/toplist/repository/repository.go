package feature_repository_toplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
)

type NatsRepository struct {
	storege *core_domain.Storage
}

func NewNatsRepository(
	storeg *core_domain.Storage,
) *NatsRepository {
	return &NatsRepository{
		storege: storeg,
	}
}
