package scrap

import (
	"context"
	"fmt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetQualityReport(
	ctx context.Context,
	sslNo string,
) (QualityReportResponse, error) {

	data, err := s.repository.GetQualityReport(ctx, sslNo)
	if err != nil {
		return QualityReportResponse{}, err
	}

	header := data.Header

	response := QualityReportResponse{
		Success:          true,
		SslSno:           header.SslSno,
		SupplierName:     header.SupplierName,
		VendorCode:       header.VendorCode,
		BillNo:           header.BillNo,
		BillDate:         header.BillDate,
		VehicleNumber:    header.VehicleNumber,
		MaterialName:     header.MaterialName,
		PartyWeight:      header.PartyWeight,
		SslWeight:        header.SslWeight,
		SslFinalWeight:   header.SslFinalWeight,
		DifferenceWeight: header.DifferenceWeight,
		TotalBagsWeight:  header.TotalBagsWeight,
		BillingRate:      header.BillingRate,
		AgentName:        header.AgentName,
		Miti:             header.Miti,
		Details:          make([]QualityReportDetailRow, 0, len(data.Details)),
		PayableDetails:   make([]QualityReportDetailRow, 0, len(data.Details)),
	}

	var payableWeight float64

	for _, d := range data.Details {
		row := QualityReportDetailRow{
			ScrapType:  d.ScrapType,
			Percentage: d.Percentage,
			Qty:        d.Qty,
			Rate:       d.Rate,
			Amount:     d.Amount,
		}

		response.Details = append(response.Details, row)

		if d.Rate != nil && *d.Rate > 0 {
			response.PayableDetails = append(response.PayableDetails, row)

			if d.Qty != nil {
				payableWeight += *d.Qty
			}
		}
	}

	response.PayableWeight = fmt.Sprintf("%.2f", payableWeight)

	return response, nil
}
