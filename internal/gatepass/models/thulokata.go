package models

type Thulokata struct {
	Sn          int64    `gorm:"column:sn;primaryKey;autoIncrement"`
	DocumentNo  string   `gorm:"column:document_no;unique;notNull"`
	PartyName   string   `gorm:"column:party_name;notNull"`
	Materials   string   `gorm:"column:materials;notNull"`
	Remarks     *string  `gorm:"column:remarks;default:-"`
	GrossWeight *float64 `gorm:"column:gross_weight"`
	TareWeight  *float64 `gorm:"column:tare_weight"`
	NetWeight   *float64 `gorm:"column:net_weight"`
	TareDate    *string  `gorm:"column:tare_date"`
	TareTime    *string  `gorm:"column:tare_time"`
	GrossDate   *string  `gorm:"column:gross_date"`
	GrossTime   *string  `gorm:"column:gross_time"`
	PrintedBy   *string  `gorm:"column:printed_by"`
}

func (Thulokata) TableName() string {
	return "thulokata"
}

type ThulokataLoading struct {
	Sn               int64    `gorm:"column:sn;primaryKey;autoIncrement"`
	ThuloKataEntryNo int64    `gorm:"column:th_entry_no;notNull"`
	Materials        string   `gorm:"column:materials;notNull"`
	InitialWeight    *float64 `gorm:"column:initial_weight"`
	FinalWeight      *float64 `gorm:"column:final_weight"`
	NetWeight        *float64 `gorm:"column:net_weight"`
	// records from the traditional thulokata_sales_order table from the legacy system
	// dont use this for new records but only for fetching the old records
	OgFinalWeight *float64 `gorm:"column:og_final_weight"`
	Bundles       int      `gorm:"column:bundles;default:0"`
	InitialDate   *string  `gorm:"column:initial_date"`
	InitialTime   *string  `gorm:"column:initial_time"`
	FinalDate     *string  `gorm:"column:final_date"`
	FinalTime     *string  `gorm:"column:final_time"`
	UserId        int      `gorm:"column:user_id;notNull"`
}

func (ThulokataLoading) TableName() string {
	return "thulokata_sales_order"
}
