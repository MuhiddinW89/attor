package analytics

import "context"

type Repository interface {
	Get(
		ctx context.Context,
		filters Filters,
	) (*Analytics, error)
}