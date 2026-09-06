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

type ThulokataRow struct {
	Materials     string
	InitialWeight *float64
	FinalWeight   *float64
	NetWeight     *float64
	Bundles       int

	PartyName string
	EntryNo   string

	ShippingAddress *string
	SalesOrder      *string
	VehicleNumber   *string
}

type SanokataRow struct {
	DocumentNo string
	Code       string
	SizeMM     string
	Bundles    int
	Pieces     int
	NetWeight  float64
	Lot        string

	LoadingDate *string

	KataNo *int

	VehicleNumber   *string
	SalesOrder      *string
	PartyName       *string
	ShippingAddress *string
}

type SHookLine struct {
	DocumentNo string
	Size       string
	Quantity   int
}
