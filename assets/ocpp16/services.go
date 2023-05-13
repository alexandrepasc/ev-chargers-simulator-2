package ocpp16

import (
	"strconv"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

/*
Set the starting configurations for the asset.
*/
func (o *Ocpp16) setConfigurations() {
	o.Conf = config

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "setConfigurations", "simulator": o.Asset.Name}, "Set startup configurations", assets.Info)

	var conTime = strconv.FormatInt(o.Timeout, 10)
	o.Conf["ConnectionTimeOut"] = core.ConfigurationKey{
		Key:      "ConnectionTimeOut",
		Readonly: true,
		Value:    &conTime,
	}
}

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

/*
Get and return the configurations set to the asset. It receives the list of configurations
requested by the central system.

Returns the list of configurations ([]core.ConfigurationKey) in case the request is an empty list
it returns all the configurations. If the request has a list it returns the requested configurations.

k	-	The requested list of configurations ([]string)
*/
func (o *Ocpp16) getConfigurationKeys(k []string) (c []core.ConfigurationKey) {
	if len(k) == 0 {
		for _, v := range config {
			c = append(c, v)
		}

		return c
	}

	for i := range k {
		_, ok := config[k[i]]

		if ok {
			c = append(c, config[k[i]])
		}
	}

	return c
}
