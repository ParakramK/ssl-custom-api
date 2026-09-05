package sales

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	rows []SalesOrderVsBillingRow
	err  error
}

func (f *fakeRepository) GetTopOutStandingCustomers(
	ctx context.Context,
	companyDB string,
	limit int,
) ([]SalesOrderVsBillingRow, error) {
	return f.rows, f.err
}

func TestGetTopOutStandingCustomers(t *testing.T) {
	repo := &fakeRepository{
		rows: []SalesOrderVsBillingRow{
			{
				CustomerCode:     "C001",
				CustomerName:     "Acme",
				OrderedQuantity:  1000,
				BilledQuantity:   1001,
				QuantityVariance: 1,
				OrderedAmount:    5000,
				BilledAmount:     5005,
				AmountVariance:   5,
			},
		},
	}

	service := NewService(repo)

	response, err := service.GetTopOutStandingCustomers(context.Background(), "DB1", 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(response.Details) != 1 {
		t.Fatalf("expected 1 row, got %d", len(response.Details))
	}

	if response.Details[0].BilledQuantity != 1001 {
		t.Errorf("expected billed quantity 1001, got %v", response.Details[0].BilledQuantity)
	}
}

func TestGetTopOutStandingCustomersEmpty(t *testing.T) {
	service := NewService(&fakeRepository{})

	response, err := service.GetTopOutStandingCustomers(context.Background(), "DB1", 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Details != nil {
		t.Errorf("expected nil details, got %+v", response.Details)
	}
}

func TestGetTopOutStandingCustomersPropagatesError(t *testing.T) {
	service := NewService(&fakeRepository{err: errors.New("db down")})

	_, err := service.GetTopOutStandingCustomers(context.Background(), "DB1", 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
