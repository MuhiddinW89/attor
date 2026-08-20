package analytics

import (
	"context"

	"github.com/MuhiddinW89/attor/internal/analytics"
)

type GetAnalyticsUseCase struct {
	service analytics.Service
}

func NewGetAnalyticsUseCase(
	service analytics.Service,
) *GetAnalyticsUseCase {
	return &GetAnalyticsUseCase{
		service: service,
	}
}

func (uc *GetAnalyticsUseCase) Execute(
	ctx context.Context,
	filters analytics.Filters,
) (*analytics.Analytics, error) {

	return uc.service.Get(
		ctx,
		filters,
	)
}