package sales

type TagDataInput struct {
	BundleTag string `query:"tag_no" doc:"Tag No Printed on the bundle tag"`
}

type TagDataOutput struct {
	Body BundleDetails `json:"body"`
}
type PackingListInput struct {
	SslSno string `query:"ssl_no" doc:"SSL No Printed on the bundle tag"`
}

type PackingListOutput struct {
	Body PackingList `json:"body"`
}

// PackingList groups the three document-number-keyed line sets:
// thulokata loadings (bundles together), sanokata loadings
// (bundles one by one) and shook lines.
type PackingList struct {
	ThulokataLines []ThulokataPackingRow `json:"thulokata_lines"`
	SanokataLines  []SanokataPackingRow  `json:"sanokata_lines"`
	ShookLines     []ShookRow            `json:"shook_lines"`
}

type ThulokataPackingRow struct {
	PartyName       string   `json:"party_name"`
	ShippingAddress string   `json:"shipping_address"`
	SalesOrder      string   `json:"sales_order"`
	VehicleNumber   string   `json:"vehicle_number"`
	EntryNo         string   `json:"entry_no"`
	Materials       string   `json:"materials"`
	InitialWeight   *float64 `json:"initial_weight"`
	FinalWeight     *float64 `json:"final_weight"`
	NetWeight       *float64 `json:"net_weight"`
	Bundles         int      `json:"bundles"`
}

type SanokataPackingRow struct {
	SalesOrder      string  `json:"sales_order"`
	PartyName       string  `json:"party_name"`
	ShippingAddress string  `json:"shipping_address"`
	VehicleNumber   string  `json:"vehicle_number"`
	OrderDate       *string `json:"order_date"`
	Code            string  `json:"code"`
	SizeMM          string  `json:"size_mm"`
	Bundles         int     `json:"bundles"`
	Pieces          int     `json:"pieces"`
	NetWeight       float64 `json:"net_weight"`
	Uid             *int64  `json:"uid"`
	KataNo          *int    `json:"kata_no"`
}

type ShookRow struct {
	Id         int64  `json:"id"`
	DocumentNo string `json:"document_no"`
	TokenNo    string `json:"token_no"`
	Size       string `json:"size"`
	Quantity   int    `json:"quantity"`
}

type BundleDetails struct {
	BundleTag       string  `json:"bundle_tag"`
	SizeMM          string  `json:"size_mm"`
	Lot             string  `json:"lot"`
	Weight          float64 `json:"weight"`
	Pieces          int     `json:"pieces"`
	Bundles         int     `json:"bundles"`
	LoadingDate     *string `json:"loading_date"`
	SslSno          string  `json:"ssl_sno"`
	VehicleNumber   string  `json:"vehicle_number"`
	LoadingSlipNo   string  `json:"loading_slip_no"`
	SalesOrder      string  `json:"sales_order"`
	PartyName       string  `json:"party_name"`
	ShippingAddress string  `json:"shipping_address"`
}
