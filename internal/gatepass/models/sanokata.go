package models

type Sanokata struct {
	Id              int64   `gorm:"column:id;primaryKey;autoIncrement"`
	DocumentNo      string  `gorm:"column:internal_numbering"`
	ShippingAddress string  `gorm:"column:shipping_address"`
	LoadingDate     *string `gorm:"column:order_date"`
	Code            string  `gorm:"column:code"`
	Lot             string  `gorm:"column:lot"`
	SizeMM          string  `gorm:"column:size_mm"`
	Bundles         int     `gorm:"column:bundles"`
	Pieces          int     `gorm:"column:pieces"`
	NetWeight       float64 `gorm:"column:net_weight"`
	CreatedAt       *string `gorm:"column:created_at;default:current_timestamp"`
	UserId          *int    `gorm:"column:user_id"`
	TokenNo         string  `gorm:"column:token_no"`

	// Belongs to gate entry: sanokata.internal_numbering -> gate_entry.document_number
	GateEntry *GateEntry `gorm:"foreignKey:DocumentNo;references:DocumentNo"`
	// Belongs to user: sales_orders.user_id -> users.id (nullable)
	User *User `gorm:"foreignKey:UserId;references:Id"`
}

func (Sanokata) TableName() string {
	return "sales_orders"
}

type SHook struct {
	Id         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	DocumentNo string `gorm:"column:document_number"`
	TokenNo    string `gorm:"column:token_no"`
	Size       string `gorm:"column:size"`
	Quantity   int    `gorm:"column:quantity"`

	// Belongs to gate entry: shook.document_number -> gate_entry.document_number
	GateEntry *GateEntry `gorm:"foreignKey:DocumentNo;references:DocumentNo"`
}

func (SHook) TableName() string {
	return "shook"
}
