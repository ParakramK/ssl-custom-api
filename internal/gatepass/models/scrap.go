package models

type Scrap struct {
	Sn               int64    `gorm:"column:sn;primaryKey;autoIncrement"`
	DocumentNo       string   `gorm:"column:document_no"`
	VehicleNumber    *string  `gorm:"column:vehicle_number"`
	BillNo           *string  `gorm:"column:bill_no"`
	BillDate         *string  `gorm:"column:bill_date"`
	EntryDate        string   `gorm:"column:entry_date"`
	EntryTime        string   `gorm:"column:entry_time"`
	PartyName        string   `gorm:"column:party_name"`
	MaterialName     string   `gorm:"column:material_name"`
	PartyWeight      float64  `gorm:"column:party_weight"`
	SslWeight        float64  `gorm:"column:ssl_weight"`
	SslFinalWeight   *float64 `gorm:"column:ssl_final_weight"`
	DifferenceWeight *float64 `gorm:"column:difference_weight"`
	SslFinalDate     *string  `gorm:"column:ssl_final_date"`
	SslFinalTime     *string  `gorm:"column:ssl_final_time"`
	SslPrintedBy     *string  `gorm:"column:ssl_printed_by"`
	CreatedAt        *string  `gorm:"column:created_at;default:current_timestamp"`
	UpdatedAt        *string  `gorm:"column:updated_at;default:current_timestamp"`
}

func (Scrap) TableName() string {
	return "scrap_entry"
}

type ScrapMultiUnload struct {
	Sn            int64    `gorm:"column:sn;primaryKey;autoIncrement"`
	ScrapSlipNo   string   `gorm:"column:scrap_slip_no;unique;notNull"`
	ScrapEntryNo  int64    `gorm:"column:scrap_entry_no"`
	InitialDate   string   `gorm:"column:initial_date"`
	InitialTime   string   `gorm:"column:initial_time"`
	InitialWeight float64  `gorm:"column:initial_weight"`
	FinalWeight   *float64 `gorm:"column:final_weight"`
	NetWeight     *float64 `gorm:"column:net_weight"`
	FinalDate     *string  `gorm:"column:final_date"`
	FinalTime     *string  `gorm:"column:final_time"`
	CreatedAt     *string  `gorm:"column:created_at;default:current_timestamp"`
	UpdatedAt     *string  `gorm:"column:updated_at;default:current_timestamp"`
}

func (ScrapMultiUnload) TableName() string {
	return "scrap_multi_unload"
}

type ScrapType struct {
	Id            int64   `gorm:"column:id;primaryKey;autoIncrement"`
	Type          string  `gorm:"column:type"`
	Rate          float64 `gorm:"column:rate"`
	UpdatedDate   *string `gorm:"column:updated_date;default:current_timestamp"`
	IncrementRate float64 `gorm:"column:increment_rate"`
}

func (ScrapType) TableName() string {
	return "scrap_type"
}
