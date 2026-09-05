package models

type Vendors struct {
	Id         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	PanNo      string `gorm:"column:pan_number"`
	VendorName string `gorm:"column:supplier_name"`
	VendorCode string `gorm:"column:CardCode;unique"`
	Active     *bool  `gorm:"column:active;default:true"`

	// Reverse of business-key reference: scrap_entry.party_name -> pan_supplier_mapping.supplier_name
	Scraps []Scrap `gorm:"foreignKey:PartyName;references:VendorName"`
}

func (Vendors) TableName() string { return "pan_supplier_mapping" }
