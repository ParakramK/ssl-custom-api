package scrap

import (
	"context"
	"errors"
	"fmt"
	"ssl-custom-api/internal/gatepass/query"
	"ssl-custom-api/internal/utils"

	"gorm.io/gorm"
)

var ErrQualityReportNotFound = errors.New("quality report not found")

type Repository interface {
	GetQualityReport(ctx context.Context, sslNo string) (*QualityReportData, error)
}

type scrapRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &scrapRepository{db: db}
}

// func (r *scrapRepository) GetQualityReport(
// 	ctx context.Context,
// 	sslNo string,
// ) (*QualityReportData, error) {

// 	if sslNo == "" {
// 		return nil, errors.New("ssl number is required")
// 	}

// 	var report models.QualityReport
// 	if err := r.db.WithContext(ctx).
// 		Model(&models.QualityReport{}).
// 		Joins("GateEntry").
// 		Joins("GateEntry.Scrap").
// 		Joins("GateEntry.Scrap.Vendor").
// 		Preload("Details", func(db *gorm.DB) *gorm.DB {
// 			return db.Order("id ASC")
// 		}).
// 		Where(&models.QualityReport{SslSno: sslNo}).
// 		Order("quality_report.created_at DESC").
// 		First(&report).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrQualityReportNotFound
// 		}
// 		return nil, fmt.Errorf("get quality report: %w", err)
// 	}

// 	// The QR -> GateEntry -> Scrap chain is required (JOINs are LEFT
// 	// JOINs, so a broken link surfaces here as nil, not as no rows).
// 	if report.GateEntry == nil || report.GateEntry.Scrap == nil {
// 		return nil, ErrQualityReportNotFound
// 	}
// 	se := report.GateEntry.Scrap
// 	ge := report.GateEntry

// 	details := report.Details
// 	if details == nil {
// 		details = []models.QualityReportDetails{}
// 	}

// 	return &QualityReportData{
// 		Header:  headers(report, se, ge),
// 		Details: details,
// 	}, nil
// }

// // headers keeps the header mapping next to its query.
// func headers(report models.QualityReport, se *models.Scrap, ge *models.GateEntry) QualityReportHeader {
// 	return QualityReportHeader{
// 		SslSno:           report.SslSno,
// 		BillNo:           se.BillNo,
// 		DriverName:       ge.DriverName,
// 		BillDate:         se.BillDate,
// 		VehicleNumber:    se.VehicleNumber,
// 		MaterialName:     se.MaterialName,
// 		PartyWeight:      se.PartyWeight,
// 		SslWeight:        se.SslWeight,
// 		SslFinalWeight:   se.SslFinalWeight,
// 		DifferenceWeight: se.DifferenceWeight,
// 		TotalBagsWeight:  report.TotalBagsWeight,
// 		BillingRate:      report.BillingRate,
// 		AgentName:        report.AgentName,
// 		Miti:             report.Miti,
// 		SupplierName:     vendorName(se),
// 		VendorCode:       vendorCode(se),
// 	}
// }

// func vendorName(se *models.Scrap) string {
// 	if se.Vendor != nil {
// 		return se.Vendor.VendorName
// 	}
// 	return ""
// }

//	func vendorCode(se *models.Scrap) string {
//		if se.Vendor != nil {
//			return se.Vendor.VendorCode
//		}
//		return ""
//	}
func (r *scrapRepository) GetQualityReport(
	ctx context.Context,
	sslNo string,
) (*QualityReportData, error) {

	if sslNo == "" {
		return nil, errors.New("ssl number is required")
	}

	q := query.Use(r.db)

	var row QualityReportResult

	err := q.QualityReport.
		WithContext(ctx).
		Select(
			// Quality report
			q.QualityReport.SslSno,
			q.QualityReport.Miti,
			q.QualityReport.AgentName,
			q.QualityReport.BillingRate,
			q.QualityReport.TotalBagsWeight,

			// Gate entry
			q.GateEntry.DriverName,

			// Scrap
			q.Scrap.BillNo,
			q.Scrap.BillDate,
			q.Scrap.VehicleNumber,
			q.Scrap.MaterialName,
			q.Scrap.PartyWeight,
			q.Scrap.SslWeight,
			q.Scrap.SslFinalWeight,
			q.Scrap.DifferenceWeight,

			// Vendor
			q.Vendors.VendorName.As("VendorName"),
			q.Vendors.VendorCode.As("VendorCode"),
		).
		Join(
			q.GateEntry,
			q.QualityReport.SslSno.EqCol(
				q.GateEntry.DocumentNo,
			),
		).
		Join(
			q.Scrap,
			q.GateEntry.DocumentNo.EqCol(
				q.Scrap.DocumentNo,
			),
		).
		Join(
			q.Vendors,
			q.Scrap.PartyName.EqCol(
				q.Vendors.VendorName,
			),
		).
		Where(
			q.QualityReport.SslSno.Eq(sslNo),
		).
		Order(
			q.QualityReport.CreatedAt.Desc(),
		).
		Limit(1).
		Scan(&row)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQualityReportNotFound
		}

		return nil, fmt.Errorf("get quality report: %w", err)
	}

	var details []QualityReportDetailResult

	err = q.QualityReportDetails.
		WithContext(ctx).
		Select(
			q.QualityReportDetails.SslSno,
			q.QualityReportDetails.ScrapType,
			q.QualityReportDetails.Percentage,
			q.QualityReportDetails.Qty,
			q.QualityReportDetails.Rate,
			q.QualityReportDetails.Amount,
		).
		Where(
			q.QualityReportDetails.SslSno.Eq(sslNo),
		).
		Order(
			q.QualityReportDetails.ID.Asc(),
		).
		Scan(&details)

	if err != nil {
		return nil, fmt.Errorf("get quality report details: %w", err)
	}

	header := QualityReportHeader{
		SslSno:           row.SslSno,
		BillNo:           row.BillNo,
		DriverName:       utils.StrVal(row.DriverName),
		BillDate:         row.BillDate,
		VehicleNumber:    row.VehicleNumber,
		MaterialName:     utils.StrVal(row.MaterialName),
		PartyWeight:      utils.ValueOrZero(row.PartyWeight),
		SslWeight:        utils.ValueOrZero(row.SslWeight),
		SslFinalWeight:   row.SslFinalWeight,
		DifferenceWeight: row.DifferenceWeight,
		TotalBagsWeight:  row.TotalBagsWeight,
		BillingRate:      row.BillingRate,
		AgentName:        row.AgentName,
		Miti:             row.Miti,
		SupplierName:     utils.StrVal(row.VendorName),
		VendorCode:       utils.StrVal(row.VendorCode),
	}

	return &QualityReportData{
		Header:  header,
		Details: details,
	}, nil
}
