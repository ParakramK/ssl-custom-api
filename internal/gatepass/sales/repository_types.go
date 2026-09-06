package sales

type TagDataRow struct {
	Code       string
	DocumentNo string
	Lot        string
	SizeMM     string
	Bundles    int
	Pieces     int
	NetWeight  float64

	GateEntryDocNo  string
	LoadingSlipNo   string
	VehicleNo       string
	SalesOrder      *string
	PartyName       *string
	ShippingAddress *string

	LoadingDate *string
}

type SHookRow struct {
	DocumentNo string
	Size       string
	Quantity   int
}
