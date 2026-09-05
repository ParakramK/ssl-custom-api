package customer

import (
	"context"
	"time"

	"ssl-custom-api/internal/utils"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetCustomerAging(
	ctx context.Context,
	cardCode string,
	companyDB string,
) (CustomerAgingDetail, error) {

	rows, err := s.repository.GetCustomerAging(
		ctx,
		cardCode,
		companyDB,
	)
	if err != nil {
		return CustomerAgingDetail{}, err
	}

	return BuildAgingResponse(rows, companyDB, false), nil
}

func BuildAgingResponse(
	rows []AgingRow,
	companyDB string,
	cached bool,
) CustomerAgingDetail {

	response := CustomerAgingDetail{
		Status:    "Success",
		CompanyDB: companyDB,
		Cached:    cached,
	}

	if len(rows) == 0 {
		return response
	}

	response.CardCode = rows[0].CustomerCode
	response.CardName = rows[0].CustomerName
	response.Balance = utils.Round2(rows[0].Balance)
	response.TaxNumber = rows[0].TaxNumber
	response.SalesEmployee = rows[0].SalesEmployee
	response.LastFetched = time.Now()

	response.Details = make([]AgingDetail, 0, len(rows))

	for _, row := range rows {
		response.Details = append(
			response.Details,
			AgingDetail{
				InvoiceNumber:     row.InvoiceNumber,
				InvoiceDate:       row.InvoiceDate,
				InvoiceAmount:     row.InvoiceAmount,
				DueDate:           row.DueDate,
				PaymentTerms:      row.PaymentTerms,
				PaymentTermGroup:  row.PaymentTermGroup,
				OutstandingDays:   row.OutstandingDays,
				PaidAmount:        row.PaidAmount,
				OutstandingAmount: row.OutstandingAmount,
				Z0To30Days:        row.Z0To30Days,
				Z31To60Days:       row.Z31To60Days,
				Z61To90Days:       row.Z61To90Days,
				Z91PlusDays:       row.Z91PlusDays,
			},
		)
	}

	response.AgingSummary = ComputeAgingSummary(rows)
	response.PaymentSummary = ComputePaymentSummary(rows)
	response.Total = ComputeTotal(rows)

	return response
}

func ComputeAgingSummary(rows []AgingRow) AgingSummary {
	var summary AgingSummary

	for _, row := range rows {
		summary.Z0To30Days += row.Z0To30Days
		summary.Z31To60Days += row.Z31To60Days
		summary.Z61To90Days += row.Z61To90Days
		summary.Z91PlusDays += row.Z91PlusDays
	}

	summary.Z0To30Days = utils.Round2(summary.Z0To30Days)
	summary.Z31To60Days = utils.Round2(summary.Z31To60Days)
	summary.Z61To90Days = utils.Round2(summary.Z61To90Days)
	summary.Z91PlusDays = utils.Round2(summary.Z91PlusDays)

	return summary
}

func ComputePaymentSummary(rows []AgingRow) PaymentSummary {
	var summary PaymentSummary

	for _, row := range rows {
		switch row.PaymentTermGroup {
		case "BG":
			summary.BG += row.OutstandingAmount
		case "LC":
			summary.LC += row.OutstandingAmount
		default:
			summary.Other += row.OutstandingAmount
		}
	}

	summary.BG = utils.Round2(summary.BG)
	summary.LC = utils.Round2(summary.LC)
	summary.Other = utils.Round2(summary.Other)

	return summary
}

func ComputeTotal(rows []AgingRow) TotalSummary {
	if len(rows) == 0 {
		return TotalSummary{}
	}

	var totalOutstanding float64

	for _, row := range rows {
		totalOutstanding += row.OutstandingAmount
	}

	accountBalance := rows[0].Balance

	return TotalSummary{
		AccountBalance:      utils.Round2(accountBalance),
		UnreconciledBalance: utils.Round2(accountBalance - totalOutstanding),
	}
}
