package sales

type TagDataInput struct {
	BundleTag string `query:"tag_no" required:"true" doc:"Tag No Printed on the bundle tag"`
}

type TagDataOutput struct {
	Body BundleDetails `json:"body"`
}
type PackingListInput struct {
	SslSno string `query:"ssl_no" required:"true" doc:"SSL Entry No"`
}

type PackingListOutput struct {
	Body PackingList `json:"body"`
}

type PackingList struct {
	Header          PackingListHeader     `json:"header"`
	ThulokataLines  []ThulokataPackingRow `json:"thulokata_lines"`
	SanokataLines   []SanokataPackingRow  `json:"sanokata_lines"`
	ShookLines      []SHookRow            `json:"shook_lines"`
	SanokataSummary []SanokataSummaryRow  `json:"sanokata_summary"`
}

type PackingListHeader struct {
	EntryNo       string `json:"entry_no"`
	PartyName     string `json:"party_name"`
	SalesOrder    string `json:"sales_order"`
	VehicleNumber string `json:"vehicle_number"`
}
type ThulokataPackingRow struct {
	Materials     string   `json:"materials"`
	InitialWeight *float64 `json:"initial_weight"`
	FinalWeight   *float64 `json:"final_weight"`
	NetWeight     *float64 `json:"net_weight"`
	Bundles       int      `json:"bundles"`
}

type SanokataPackingRow struct {
	LoadingDate *string `json:"loading_date"`
	Code        string  `json:"code"`
	SizeMM      string  `json:"size_mm"`
	Bundles     int     `json:"bundles"`
	Pieces      int     `json:"pieces"`
	NetWeight   float64 `json:"net_weight"`
	KataNo      *int    `json:"kata_no"`
}
type SanokataSummaryRow struct {
	SizeMM    string  `json:"size_mm"`
	Bundles   int     `json:"bundles"`
	NetWeight float64 `json:"net_weight"`
}

type SHookRow struct {
	DocumentNo string `json:"document_no"`
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
