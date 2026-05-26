package feature_service_stoplist

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	"context"
	"fmt"
)

func (s *ServiceStopList) GetStopList(
	ctx context.Context,
) (*core_domain.StopList, error) {
	list := s.stopList.Items()
	reqStopList := core_domain.NewStopList(list)
	fmt.Println(reqStopList, "in service")

	return reqStopList, nil
}
