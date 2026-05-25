package feature_service_toplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"context"
)

type NatsRepository interface {
	GetList(
		ctx context.Context,
		limit int,
	) ([]core_domain.TopItem, error)
}

type NatsService struct {
	repository NatsRepository
}

func NewNatsService(
	repository NatsRepository,
) *NatsService {
	return &NatsService{
		repository: repository,
	}
}
