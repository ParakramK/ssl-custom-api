package scrap

type QualityReportData struct {
	Header  QualityReportHeader
	Details []QualityReportDetailResult
}
type QualityReportResult struct {
	SslSno          string
	Miti            *string
	AgentName       *string
	BillingRate     *float64
	TotalBagsWeight *float64

	BillNo           *string
	BillDate         *string
	VehicleNumber    *string
	MaterialName     *string
	PartyWeight      *float64
	SslWeight        *float64
	SslFinalWeight   *float64
	DifferenceWeight *float64

	DriverName *string

	VendorName *string
	VendorCode *string
}

type QualityReportDetailResult struct {
	SslSno     *string
	ScrapType  *string
	Percentage *float64
	Qty        *float64
	Rate       *float64
	Amount     *float64
}
