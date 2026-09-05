package sales

import "context"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetTopOutStandingCustomers(
	ctx context.Context,
	companyDB string,
	limit int,
) (TopCustomerResponse, error) {

	rows, err := s.repository.GetTopOutStandingCustomers(
		ctx,
		companyDB,
		limit,
	)
	if err != nil {
		return TopCustomerResponse{}, err
	}

	return BuildResponse(rows, companyDB, false), nil
}

func BuildResponse(
	rows []SalesOrderVsBillingRow,
	companyDB string,
	cached bool,
) TopCustomerResponse {

	return TopCustomerResponse{
		Details: rows,
		Cached:  cached,
	}
}
