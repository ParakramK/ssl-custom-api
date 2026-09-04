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
}
