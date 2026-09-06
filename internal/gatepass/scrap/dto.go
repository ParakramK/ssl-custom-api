package scrap

type QualityReportInput struct {
	DocumentNo string `query:"ssl_no" doc:"Gate entry document number"`
}

type QualityReportOutput struct {
	Body QualityReportResponse
}

type QualityReportResponse struct {
	Success          bool                     `json:"success"`
	SslSno           string                   `json:"ssl_sno"`
	SupplierName     string                   `json:"supplier_name"`
	VendorCode       string                   `json:"vendor_code"`
	BillNo           *string                  `json:"bill_no"`
	BillDate         *string                  `json:"bill_date"`
	VehicleNumber    *string                  `json:"vehicle_number"`
	MaterialName     string                   `json:"material_name"`
	PartyWeight      float64                  `json:"party_weight"`
	SslWeight        float64                  `json:"ssl_weight"`
	SslFinalWeight   *float64                 `json:"ssl_final_weight"`
	DifferenceWeight *float64                 `json:"difference_weight"`
	TotalBagsWeight  *float64                 `json:"total_bags_weight"`
	BillingRate      *float64                 `json:"billing_rate"`
	AgentName        *string                  `json:"agent_name"`
	Miti             *string                  `json:"miti"`
	PayableWeight    string                   `json:"payable_weight"`
	Details          []QualityReportDetailRow `json:"details"`
	PayableDetails   []QualityReportDetailRow `json:"payable_details"`
}

type QualityReportDetailRow struct {
	ID         int64    `json:"id"`
	ScrapType  *string  `json:"scrap_type"`
	Percentage *float64 `json:"percentage"`
	Qty        *float64 `json:"qty"`
	Rate       *float64 `json:"rate"`
	Amount     *float64 `json:"amount"`
}

type QualityReportHeader struct {
	SslSno           string   `gorm:"column:ssl_sno"`
	DriverName       string   `gorm:"column:driver_name"`
	SupplierName     string   `gorm:"column:supplier_name"`
	VendorCode       string   `gorm:"column:vendor_code"`
	BillNo           *string  `gorm:"column:bill_no"`
	BillDate         *string  `gorm:"column:bill_date"`
	VehicleNumber    *string  `gorm:"column:vehicle_number"`
	MaterialName     string   `gorm:"column:material_name"`
	PartyWeight      float64  `gorm:"column:party_weight"`
	SslWeight        float64  `gorm:"column:ssl_weight"`
	SslFinalWeight   *float64 `gorm:"column:ssl_final_weight"`
	DifferenceWeight *float64 `gorm:"column:difference_weight"`
	TotalBagsWeight  *float64 `gorm:"column:total_bags_weight"`
	BillingRate      *float64 `gorm:"column:billing_rate"`
	AgentName        *string  `gorm:"column:agent_name"`
	Miti             *string  `gorm:"column:miti"`
}
