package sales

import (
	"context"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetTagData(
	ctx context.Context,
	sslNo string,
) (BundleDetails, error) {

	data, err := s.repository.GetTagData(ctx, sslNo)
	if err != nil {
		return BundleDetails{}, err
	}
	if data == nil {
		return BundleDetails{}, ErrTagNotFound
	}

	return *data, nil
}

func (s *Service) GetPackingList(
	ctx context.Context,
	sslNo string,
) (*PackingList, error) {

	return s.repository.GetPackingList(ctx, sslNo)
}
