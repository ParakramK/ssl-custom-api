package models

type QualityReport struct {
	SslSno          string   `gorm:"column:ssl_sno;primaryKey;unique;notNull"`
	QualityId       *string  `gorm:"column:quality_id;default:-"`
	TokenNo         *string  `gorm:"column:token_no"`
	Miti            *string  `gorm:"column:miti"`
	AgentName       *string  `gorm:"column:agent_name"`
	BillingRate     *float64 `gorm:"column:billing_rate"`
	EndCutDetails   *string  `gorm:"column:end_cut_details"`
	OverSizeDetails *string  `gorm:"column:over_size_details"`
	Bags            *string  `gorm:"column:bags"`
	TotalBagsWeight *float64 `gorm:"column:total_bags_weight"`

	CreatedAt   *string `gorm:"column:created_at;default:current_timestamp"`
	Updated_at  *string `gorm:"column:updated_at;default:current_timestamp"`
	PoCreated   *bool   `gorm:"column:po_created;default:false"`
	PostingDate *string `gorm:"column:posting_date"`
}

func (QualityReport) TableName() string {
	return "quality_report"
}

type QualityReportDetails struct {
	ID         int64    `gorm:"column:id;primaryKey;autoIncrement"`
	SslSno     *string  `gorm:"column:ssl_sno"`
	ScrapType  *string  `gorm:"column:scrap_type"`
	Percentage *float64 `gorm:"column:percentage"`
	Qty        *float64 `gorm:"column:qty"`
	Rate       *float64 `gorm:"column:rate"`
	Amount     *float64 `gorm:"column:amount"`
}

func (QualityReportDetails) TableName() string {
	return "quality_report_details"
}
