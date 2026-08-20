package analytics

import "context"

type Service interface {
	Get(
		ctx context.Context,
		filters Filters,
	) (*Analytics, error)
}