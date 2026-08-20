package analytics

import "context"

type service struct {
	repo Repository
}

func NewService(
	repo Repository,
) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Get(
	ctx context.Context,
	filters Filters,
) (*Analytics, error) {

	return s.repo.Get(
		ctx,
		filters,
	)
}