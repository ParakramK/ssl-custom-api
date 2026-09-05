package models

type GateEntry struct {
	DocumentNo     string  `gorm:"column:document_number;primaryKey"`
	TokenId        *string `gorm:"column:token_id"`
	LoadingSlipNo  string  `gorm:"column:loading_slip_no;default:0"`
	VehicleNo      string  `gorm:"column:vehicle_number"`
	DriverName     string  `gorm:"column:driver_name"`
	DriverMobile   string  `gorm:"column:driver_mobile"`
	Purpose        string  `gorm:"column:purpose"`
	EntryDate      *string `gorm:"column:entry_date"`
	EntryTime      *string `gorm:"column:entry_time"`
	ExitDate       *string `gorm:"column:exit_date"`
	ExitTime       *string `gorm:"column:exit_time"`
	GateUser       *int    `gorm:"column:gate_user"`
	ThulokataEntry *int    `gorm:"column:thulokata_entry"`
	ScrapEntry     *int    `gorm:"column:scrap_entry"`

	// Business-key references via shared entry/document number
	// gate_entry.document_number -> scrap_entry.document_no
	Scrap *Scrap `gorm:"foreignKey:DocumentNo;references:DocumentNo"`
	// gate_entry.document_number -> thulokata.document_no
	Thulokata *Thulokata `gorm:"foreignKey:DocumentNo;references:DocumentNo"`
	// gate_entry.document_number -> quality_report.ssl_sno
	QualityReport *QualityReport `gorm:"foreignKey:DocumentNo;references:SslSno"`
	// Business-key reference: gate_entry.loading_slip_no -> loading_entry.token
	LoadingEntry *LoadingEntry `gorm:"foreignKey:LoadingSlipNo;references:LoadingSlipNo"`
	// One gate entry has many sanokata rows: sanokata.internal_numbering -> gate_entry.document_number
	Sanokatas []Sanokata `gorm:"foreignKey:DocumentNo;references:DocumentNo"`
	// One gate entry has many shook rows: shook.document_number -> gate_entry.document_number
	SHooks []SHook `gorm:"foreignKey:DocumentNo;references:DocumentNo"`
}
