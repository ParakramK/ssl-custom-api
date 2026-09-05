package models

type LoadingEntry struct {
	ID              int64   `gorm:"column:id;primaryKey;autoIncrement"`
	LoadingSlipNo   string  `gorm:"column:token;unique;notNull"`
	BookletNo       *string `gorm:"column:booklet_no"`
	PartyName       *string `gorm:"column:party_name"`
	VatNo           *string `gorm:"column:vat"`
	VehicleNo       *string `gorm:"column:vehicle_number"`
	SalesOrder      *string `gorm:"column:sales_order"`
	Transporter     *string `gorm:"column:transporter"`
	CreatedAt       *string `gorm:"column:created_at;default:current_timestamp"`
	ShippingAddress *string `gorm:"column:shipping_address"`
	Type            *string `gorm:"column:type"`
	BillingComplete *bool   `gorm:"column:billing_complete;default:false"`
	SoDocEntry      *int64  `gorm:"column:so_doc_entry"`

	// One loading entry has many detail lines: loading_entry_details.loading_entry_id -> loading_entry.id
	Details []LoadingEntryDetails `gorm:"foreignKey:LoadingEntryID;references:ID"`
	// Gate entries pointing at this loading slip: gate_entry.loading_slip_no -> loading_entry.token
	GateEntries []GateEntry `gorm:"foreignKey:LoadingSlipNo;references:LoadingSlipNo"`
}

type LoadingEntryDetails struct {
	ID             int64   `gorm:"column:id;primaryKey;autoIncrement"`
	LoadingEntryID int64   `gorm:"column:loading_entry_id"`
	Sn             *string `gorm:"column:sn"`
	Item           *string `gorm:"column:item"`
	Weight         *string `gorm:"column:weight"`
	Bundles        *int    `gorm:"column:bundles"`
	ItemCode       *string `gorm:"column:itemcode"`

	// Belongs to parent loading entry: loading_entry_id -> loading_entry.id
	LoadingEntry *LoadingEntry `gorm:"foreignKey:LoadingEntryID;references:ID"`
}

type GateLoadingSlip struct {
	ID            int64   `gorm:"column:id;primaryKey;autoIncrement"`
	LoadingSlipNo string  `gorm:"column:loading_slip_no"`
	PartyName     string  `gorm:"column:party_name"`
	VehicleNumber string  `gorm:"column:vehicle_number"`
	CreatedAt     *string `gorm:"column:created_at;default:current_timestamp"`
}
