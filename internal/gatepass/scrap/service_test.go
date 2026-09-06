package scrap

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	data     *QualityReportData
	err      error
	scrapErr error
}

func (f *fakeRepository) GetQualityReport(
	ctx context.Context,
	sslNo string,
) (*QualityReportData, error) {
	return f.data, f.err
}

func floatPtr(v float64) *float64 {
	return &v
}

func strPtr(v string) *string {
	return &v
}

func TestGetQualityReportPayableMath(t *testing.T) {
	repo := &fakeRepository{
		data: &QualityReportData{
			Header: QualityReportHeader{
				SslSno:       "SSL001",
				SupplierName: "Acme",
				VendorCode:   "C001",
				PartyWeight:  1000,
				SslWeight:    999,
			},
			Details: []QualityReportDetailResult{
				{ScrapType: strPtr("A"), Qty: floatPtr(100.5), Rate: floatPtr(10)},
				{ScrapType: strPtr("B"), Qty: floatPtr(50), Rate: floatPtr(0)},
				{ScrapType: strPtr("C"), Qty: floatPtr(25.25), Rate: nil},
				{ScrapType: strPtr("D"), Qty: nil, Rate: floatPtr(5)},
			},
		},
	}

	service := NewService(repo)

	response, err := service.GetQualityReport(context.Background(), "SSL001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("expected success true")
	}

	if len(response.Details) != 4 {
		t.Fatalf("expected 4 details, got %d", len(response.Details))
	}

	if len(response.PayableDetails) != 2 {
		t.Fatalf("expected 2 payable details, got %d", len(response.PayableDetails))
	}

	if response.PayableWeight != "100.50" {
		t.Errorf("expected payable weight 100.50, got %q", response.PayableWeight)
	}
}

func TestGetQualityReportEmptyDetails(t *testing.T) {
	repo := &fakeRepository{
		data: &QualityReportData{
			Header: QualityReportHeader{SslSno: "SSL002"},
		},
	}

	service := NewService(repo)

	response, err := service.GetQualityReport(context.Background(), "SSL002")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Details == nil || response.PayableDetails == nil {
		t.Error("expected empty lists, not null")
	}

	if response.PayableWeight != "0.00" {
		t.Errorf("expected payable weight 0.00, got %q", response.PayableWeight)
	}
}

func TestGetQualityReportPropagatesNotFound(t *testing.T) {
	service := NewService(&fakeRepository{err: ErrQualityReportNotFound})

	_, err := service.GetQualityReport(context.Background(), "SSL999")
	if !errors.Is(err, ErrQualityReportNotFound) {
		t.Fatalf("expected not-found error, got %v", err)
	}
}

func TestGetQualityReportPropagatesError(t *testing.T) {
	service := NewService(&fakeRepository{err: errors.New("db down")})

	_, err := service.GetQualityReport(context.Background(), "SSL001")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
