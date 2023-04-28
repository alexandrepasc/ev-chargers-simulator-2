package ocpp16

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

/*
Sends the boot notification to the central system.
*/
func (o *Ocpp16) sendBootNotification() {
	var bn = core.BootNotificationRequest{
		ChargePointSerialNumber: o.Mod.Ocpp.SerialNumb,
		ChargePointModel:        o.Mod.Ocpp.Model,
		ChargePointVendor:       o.Mod.Ocpp.Vendor,
		FirmwareVersion:         o.Mod.Ocpp.FwVersion,
		MeterSerialNumber:       o.Mod.Ocpp.MeterSerialNumb,
		Iccid:                   o.Mod.Ocpp.Modem.Iccid,
		Imsi:                    o.Mod.Ocpp.Modem.Imsi,
	}

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "sendBootNotification", "model": o.Asset.Name}, bn, assets.Info)

	resp, err := o.s.SendRequest(bn)

	if err != nil {
		o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "sendBootNotification", "model": o.Asset.Name}, err, assets.Error)

		return
	}

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "sendBootNotification", "model": o.Asset.Name}, resp, assets.Info)
}
