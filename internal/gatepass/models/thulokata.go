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

	// One thulokata has many loading rows: thulokata_sales_order.th_entry_no -> thulokata.sn
	Loadings []ThulokataLoading `gorm:"foreignKey:ThuloKataEntryNo;references:Sn"`
	// Gate entries pointing at this record via sn FK: gate_entry.thulokata_entry -> thulokata.sn
	GateEntries []GateEntry `gorm:"foreignKey:ThulokataEntry;references:Sn"`
	// Business-key reference via shared entry number: thulokata.document_no -> gate_entry.document_number
	GateEntry *GateEntry `gorm:"foreignKey:DocumentNo;references:DocumentNo"`
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

	// Belongs to parent thulokata: th_entry_no -> thulokata.sn
	Thulokata *Thulokata `gorm:"foreignKey:ThuloKataEntryNo;references:Sn"`
}

func (ThulokataLoading) TableName() string {
	return "thulokata_sales_order"
}
