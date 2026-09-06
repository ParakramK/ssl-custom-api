package sales

import (
	"ssl-custom-api/internal/gatepass/models"
	"ssl-custom-api/internal/utils"
)

func TagDetails(tag TagDataRow) BundleDetails {
	return BundleDetails{
		BundleTag:       tag.Code,
		SizeMM:          tag.SizeMM,
		Lot:             tag.Lot,
		Weight:          tag.NetWeight,
		Pieces:          tag.Pieces,
		Bundles:         tag.Bundles,
		LoadingDate:     tag.LoadingDate,
		SslSno:          tag.DocumentNo,
		LoadingSlipNo:   tag.LoadingSlipNo,
		VehicleNumber:   tag.VehicleNo,
		SalesOrder:      utils.StrVal(tag.SalesOrder),
		PartyName:       utils.StrVal(tag.PartyName),
		ShippingAddress: utils.StrVal(tag.ShippingAddress),
	}
}

func packThulokataRow(thl models.ThulokataLoading) ThulokataPackingRow {
	row := ThulokataPackingRow{
		Materials:     thl.Materials,
		InitialWeight: thl.InitialWeight,
		FinalWeight:   thl.FinalWeight,
		NetWeight:     thl.NetWeight,
		Bundles:       thl.Bundles,
	}
	if th := thl.Thulokata; th != nil {
		row.PartyName = th.PartyName
		row.EntryNo = th.DocumentNo
		if ge := th.GateEntry; ge != nil {
			if le := ge.LoadingEntry; le != nil {
				row.ShippingAddress = utils.StrVal(le.ShippingAddress)
				row.SalesOrder = utils.StrVal(le.SalesOrder)
				row.VehicleNumber = utils.StrVal(le.VehicleNo)
			}
		}
	}
	return row
}

func packSanokataRow(so models.Sanokata) SanokataPackingRow {
	row := SanokataPackingRow{
		OrderDate: so.LoadingDate,
		Code:      so.Code,
		SizeMM:    so.SizeMM,
		Bundles:   so.Bundles,
		Pieces:    so.Pieces,
		NetWeight: so.NetWeight,
	}
	if u := so.User; u != nil {
		row.Uid = &u.Id
		row.KataNo = u.KataNo
	}
	if ge := so.GateEntry; ge != nil {
		row.VehicleNumber = ge.VehicleNo
		if le := ge.LoadingEntry; le != nil {
			row.SalesOrder = utils.StrVal(le.SalesOrder)
			row.PartyName = utils.StrVal(le.PartyName)
			row.ShippingAddress = utils.StrVal(le.ShippingAddress)
		}
	}
	return row
}

func packShookRow(sh ShookRow) ShookRow {
	return ShookRow{
		DocumentNo: sh.DocumentNo,
		Size:       sh.Size,
		Quantity:   sh.Quantity,
	}
}
