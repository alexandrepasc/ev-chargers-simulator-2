package ocpp16

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

/*
Set the starting configurations for the asset.
*/
func (o *Ocpp16) setStartUpConfigurations() {
	o.Conf = config

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "setConfigurations", "simulator": o.Asset.Name}, "Set startup configurations", assets.Info)

	var t = strconv.FormatInt(o.Timeout, 10)

	var cto = o.Conf["ConnectionTimeOut"]

	cto.Value = &t

	o.Conf["ConnectionTimeOut"] = cto

	var mvi = "20"

	var mvsi = o.Conf["MeterValueSampleInterval"]

	mvsi.Value = &mvi
	o.Conf["MeterValueSampleInterval"] = mvsi

	o.t = 0
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
Process the start transaction request and proceed with the correct logic in each case of the
authorize config.

Validates if the connector selected is not active, activate it. If no connector were selected
validate if any of them can be activated, if so activate it.
*/
func (o *Ocpp16) processRemoteStartTransaction(r *core.RemoteStartTransactionRequest) *core.RemoteStartTransactionConfirmation {
	if !o.Mod.Ocpp.AuthorizeRemote {
		o.chargeProfile = r.ChargingProfile

		// At the moment not know how to identify the evse from the request so will only consider 1
		if r.ConnectorId != nil {
			var ok = true

			for _, c := range o.Asset.Evses[0].Connectors {
				if c.Enabled {
					ok = false
				}
			}

			if ok {
				for i, c := range o.Asset.Evses[0].Connectors {
					if c.ID == int64(*r.ConnectorId) {
						if !c.Enabled {
							o.Asset.Evses[0].Connectors[i].Enabled = true

							var req = core.StartTransactionRequest{
								ConnectorId: int(c.ID),
								IdTag:       r.IdTag,
								Timestamp:   types.NewDateTime(time.Now()),
							}

							go o.s.SendRequest(req) //nolint:errcheck // because at the moment can not handle the error since it is in a routine

							return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}
						}
					}
				}
			} else {
				return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}
			}
		}

		var ok = true

		for _, c := range o.Asset.Evses[0].Connectors {
			if c.Enabled {
				ok = false
			}
		}

		if ok {
			for i, c := range o.Asset.Evses[0].Connectors {
				if !c.Enabled {
					o.Asset.Evses[0].Connectors[i].Enabled = true

					var req = core.StartTransactionRequest{
						ConnectorId: int(c.ID),
						IdTag:       r.IdTag,
						Timestamp:   types.NewDateTime(time.Now()),
					}

					go o.s.SendRequest(req) //nolint:errcheck // because at the moment can not handle the error since it is in a routine

					return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}
				}
			}
		} else {
			return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}
		}
	}
	// TODO: need the logic when it needs to authenticate

	return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}
}

/**/
func (o *Ocpp16) updateData() {
	for x, e := range o.Asset.Evses {
		for y, c := range e.Connectors {
			if c.Enabled {
				fmt.Println("updateData")

				if c.DP.Ticker <= c.Data[c.DP.Position].Duration {
					o.Asset.Evses[x].Connectors[y].DP.Ticker++
				} else {
					o.Asset.Evses[x].Connectors[y].DP.Ticker = 0

					if c.DP.Position < int64(len(c.Data)-1) {
						o.Asset.Evses[x].Connectors[y].DP.Position++
					} else {
						o.Asset.Evses[x].Connectors[y].DP.Position = 0
					}
				}
			}
		}
	}
}

func (o *Ocpp16) meterValues() {
	fmt.Println("meter values")
}

/*
Get and return the configurations set to the asset. It receives the list of configurations
requested by the central system.

Returns the list of configurations ([]core.ConfigurationKey) in case the request is an empty list
it returns all the configurations. If the request has a list it returns the requested configurations.

It also returns the list of keys requested that are not configured ([]string).

k	-	The requested list of configurations ([]string)
*/
func (o *Ocpp16) getConfigurationKeys(k []string) (c []core.ConfigurationKey, u []string) {
	if len(k) == 0 {
		for _, v := range o.Conf {
			c = append(c, v)
		}

		return c, u
	}

	for i := range k {
		_, ok := o.Conf[k[i]]

		if ok {
			c = append(c, o.Conf[k[i]])
		} else {
			u = append(u, k[i])
		}
	}

	return c, u
}

/*
Validates if the configuration key is supported and if it is writable. If it passes the previous
cases will set the new value to the configurations and returns accepted (core.ConfigurationStatus).

If the key does not match the list of supported keys will return not supported.

If the key is read only will return rejected.

c	-	Configuration key and value sent from the central system (*core.ChangeConfigurationRequest).
*/
func (o *Ocpp16) setConfiguration(c *core.ChangeConfigurationRequest) core.ConfigurationStatus {
	k, ok := o.Conf[c.Key]

	if ok {
		if !k.Readonly {
			k.Value = &c.Value
			o.Conf[c.Key] = k
		} else {
			return core.ConfigurationStatusRejected
		}
	} else {
		return core.ConfigurationStatusNotSupported
	}

	return core.ConfigurationStatusAccepted
}

func (o *Ocpp16) handleTick() {
	const rInt64 = math.MaxInt64 - 7

	if o.t > rInt64 {
		o.t = 0
	} else {
		o.t++
	}
}
