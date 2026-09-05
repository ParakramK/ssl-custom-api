package customer

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	rows []AgingRow
	err  error
}

func (f *fakeRepository) GetCustomerAging(
	ctx context.Context,
	cardCode string,
	companyDB string,
) ([]AgingRow, error) {
	return f.rows, f.err
}

func TestGetCustomerAgingAggregatesRows(t *testing.T) {
	repo := &fakeRepository{
		rows: []AgingRow{
			{
				CustomerCode:      "C001",
				CustomerName:      "Acme",
				Balance:           1000.456,
				TaxNumber:         "T1",
				SalesEmployee:     "Jane",
				PaymentTermGroup:  "BG",
				OutstandingAmount: 100,
				Z0To30Days:        100,
			},
			{
				CustomerCode:      "C001",
				CustomerName:      "Acme",
				Balance:           1000.456,
				TaxNumber:         "T1",
				SalesEmployee:     "Jane",
				PaymentTermGroup:  "LC",
				OutstandingAmount: 50,
				Z31To60Days:       50,
			},
		},
	}

	service := NewService(repo)

	response, err := service.GetCustomerAging(context.Background(), "C001", "DB1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Status != "Success" {
		t.Errorf("expected Status Success, got %q", response.Status)
	}

	if response.CardCode != "C001" {
		t.Errorf("expected CardCode C001, got %q", response.CardCode)
	}

	if len(response.Details) != 2 {
		t.Fatalf("expected 2 details, got %d", len(response.Details))
	}

	if response.Balance != 1000.46 {
		t.Errorf("expected rounded balance 1000.46, got %v", response.Balance)
	}

	if response.AgingSummary.Z0To30Days != 100 {
		t.Errorf("expected z0to30 100, got %v", response.AgingSummary.Z0To30Days)
	}

	if response.PaymentSummary.BG != 100 || response.PaymentSummary.LC != 50 {
		t.Errorf("unexpected payment summary: %+v", response.PaymentSummary)
	}

	if response.Total.AccountBalance != 1000.46 {
		t.Errorf("expected account balance 1000.46, got %v", response.Total.AccountBalance)
	}

	if response.Total.UnreconciledBalance != 850.46 {
		t.Errorf("expected unreconciled 850.46, got %v", response.Total.UnreconciledBalance)
	}
}

func TestGetCustomerAgingEmpty(t *testing.T) {
	service := NewService(&fakeRepository{})

	response, err := service.GetCustomerAging(context.Background(), "C999", "DB1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Status != "Success" {
		t.Errorf("expected Status Success, got %q", response.Status)
	}

	if response.Details != nil {
		t.Errorf("expected nil details, got %+v", response.Details)
	}
}

func TestGetCustomerAgingPropagatesError(t *testing.T) {
	service := NewService(&fakeRepository{err: errors.New("db down")})

	_, err := service.GetCustomerAging(context.Background(), "C001", "DB1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestComputeTotalEmpty(t *testing.T) {
	total := ComputeTotal(nil)

	if total != (TotalSummary{}) {
		t.Errorf("expected zero total, got %+v", total)
	}
}
