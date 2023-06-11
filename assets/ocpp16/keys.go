package ocpp16

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

var config = map[string]core.ConfigurationKey{
	"ChargingScheduleAllowedChargingRateUnit": {
		Key:      "ChargingScheduleAllowedChargingRateUnit",
		Readonly: false,
		Value:    &assets.Current,
	},
	"ConnectionTimeOut": {
		Key:      "ConnectionTimeOut",
		Readonly: false,
	},
	"MeterValuesSampledData": {
		Key:      "MeterValuesSampledData",
		Readonly: false,
		Value:    &assets.EnergyActiveImportRegister,
	},
	"MeterValueSampleInterval": {
		Key:      "MeterValueSampleInterval",
		Readonly: false,
	},
}
