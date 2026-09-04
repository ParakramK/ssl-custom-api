package models

type Sanokata struct {
	Id              int64   `gorm:"column:id;primaryKey;autoIncrement"`
	DocumentNo      string  `gorm:"column:internal_numbering"`
	ShippingAddress string  `gorm:"column:shipping_address"`
	OrderDate       *string `gorm:"column:order_date"`
	Code            string  `gorm:"column:code"`
	Lot             string  `gorm:"column:lot"`
	SizeMM          string  `gorm:"column:size_mm"`
	Bundles         int     `gorm:"column:bundles"`
	Pieces          int     `gorm:"column:pieces"`
	NetWeight       float64 `gorm:"column:net_weight"`
	CreatedAt       *string `gorm:"column:created_at;default:current_timestamp"`
	UserId          *int    `gorm:"column:user_id"`
	TokenNo         string  `gorm:"column:token_no"`
}

func (Sanokata) TableName() string {
	return "sales_order"
}

type SHook struct {
	Id         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	DocumentNo string `gorm:"column:document_number"`
	TokenNo    string `gorm:"column:token_no"`
	Size       string `gorm:"column:size"`
	Quantity   int    `gorm:"column:quantity"`
}

func (SHook) TableName() string {
	return "shook"
}
