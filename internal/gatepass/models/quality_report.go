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

	// One report has many detail lines: quality_report_details.ssl_sno -> quality_report.ssl_sno
	Details []QualityReportDetails `gorm:"foreignKey:SslSno;references:SslSno"`
	// Business-key reference: quality_report.agent_name -> agent.agent (nullable, may not match)
	Agent *Agent `gorm:"foreignKey:AgentName;references:Agent"`
	// Business-key reference via shared entry number: quality_report.ssl_sno -> gate_entry.document_number
	GateEntry *GateEntry `gorm:"foreignKey:SslSno;references:DocumentNo"`
}

func (QualityReport) TableName() string { return "quality_report" }

type QualityReportDetails struct {
	ID         int64    `gorm:"column:id;primaryKey;autoIncrement"`
	SslSno     *string  `gorm:"column:ssl_sno"`
	ScrapType  *string  `gorm:"column:scrap_type"`
	Percentage *float64 `gorm:"column:percentage"`
	Qty        *float64 `gorm:"column:qty"`
	Rate       *float64 `gorm:"column:rate"`
	Amount     *float64 `gorm:"column:amount"`

	// Belongs to parent report: quality_report_details.ssl_sno -> quality_report.ssl_sno
	Report *QualityReport `gorm:"foreignKey:SslSno;references:SslSno"`
	// Business-key reference: quality_report_details.scrap_type -> scrap_types.type (nullable)
	ScrapTypeInfo *ScrapType `gorm:"foreignKey:ScrapType;references:Type"`
}

func (QualityReportDetails) TableName() string { return "quality_report_details" }
