package ocpp16

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/lorenzodonini/ocpp-go/ocpp"
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

	var mvsdc = assets.EnergyActiveImportRegister + "," + assets.PowerActiveImport

	var mvsd = o.Conf["MeterValuesSampledData"]

	mvsd.Value = &mvsdc
	o.Conf["MeterValuesSampledData"] = mvsd

	o.t = 0

	for x, e := range o.Asset.Evses {
		for y := range e.Connectors {
			o.Asset.Evses[x].Connectors[y].DP.Position = 0
			o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
			o.Asset.Evses[x].Connectors[y].TPower = 0
			o.Asset.Evses[x].Connectors[y].Energy = 0
			o.Asset.Evses[x].Connectors[y].Enabled = false
		}
	}

	o.st = time.Now()
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
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "processRemoteStartTransaction",
		"feature":   "StartTransaction",
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

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
					if c.ID != int64(*r.ConnectorId) {
						continue
					}

					o.Asset.Evses[0].Connectors[i].Enabled = true

					o.Asset.Evses[0].Connectors[i].DP.Position = 2
					o.Asset.Evses[0].Connectors[i].DP.Ticker = 0

					fmt.Println(o.Asset.Evses[0].Connectors[i].DP)

					var req = core.StartTransactionRequest{
						ConnectorId: int(c.ID),
						IdTag:       r.IdTag,
						MeterStart:  int(c.Energy),
						Timestamp:   types.NewDateTime(time.Now()),
					}

					lm["feature"] = req.GetFeatureName()

					o.logger.log(lm, req, assets.Error)

					cb := func(res ocpp.Response, err error) {
						var lm2 = map[string]string{
							"protocol":  string(o.Asset.Protocol),
							"function":  "processRemoteStartTransaction",
							"feature":   res.GetFeatureName(),
							"simulator": o.Asset.Name,
							"sender":    assets.CS,
							"type":      assets.Response,
						}

						o.logger.log(lm2, res, assets.Info)
					}
					err := o.s.SendRequestAsync(req, cb)

					lm["type"] = assets.Response

					if err != nil {
						o.logger.log(lm, err, assets.Error)
					}

					o.statusNotification(o.Asset.Evses[0].Connectors[i])

					return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}
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
				o.Asset.Evses[0].Connectors[i].Enabled = true

				o.Asset.Evses[0].Connectors[i].DP.Position = 2

				var req = core.StartTransactionRequest{
					ConnectorId: int(c.ID),
					IdTag:       r.IdTag,
					Timestamp:   types.NewDateTime(time.Now()),
				}

				go o.s.SendRequest(req) //nolint:errcheck // because at the moment can not handle the error since it is in a routine

				return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}
			}
		} else {
			return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}
		}
	}
	// TODO: need the logic when it needs to authenticate

	return &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}
}

/**/
func (o *Ocpp16) processRemoteStopTransaction(r *core.RemoteStopTransactionRequest) *core.RemoteStopTransactionConfirmation {
	if r.TransactionId == o.chargeProfile.TransactionId {
		// at the moment this is only supporting 1 evse per simulator, so this will only look for one position
		for i, c := range o.Asset.Evses[0].Connectors {
			if !c.Enabled {
				continue
			}

			o.Asset.Evses[0].Connectors[i].DP.Position = int64(len(c.Data) - 1)
			o.Asset.Evses[0].Connectors[i].DP.Ticker = 0

			go o.statusNotification(o.Asset.Evses[0].Connectors[i])

			go o.stopTransaction(c)

			return &core.RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}
		}
	}

	return &core.RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}
}

/**/
// TODO: the total power calculation need to be reviewed, at the moment with 100 w in a couple of secs the result is 0
func (o *Ocpp16) updateData() {
	for x, e := range o.Asset.Evses {
		for y, c := range e.Connectors {
			if o.Asset.StartCharging {
				if canEnable(e.Connectors) {
					o.Asset.Evses[x].Connectors[y].Enabled = true
				}

				if c.Enabled {
					var cs = c.Data[c.DP.Position].ChargingState

					if c.DP.Ticker < c.Data[c.DP.Position].Duration {
						o.Asset.Evses[x].Connectors[y].DP.Ticker++
					} else {
						o.Asset.Evses[x].Connectors[y].DP.Ticker = 0

						if c.DP.Position < int64(len(c.Data)-1) {
							o.Asset.Evses[x].Connectors[y].DP.Position++
						} else {
							o.Asset.Evses[x].Connectors[y].DP.Position = 0
						}
					}

					if cs != c.Data[c.DP.Position].ChargingState {
						o.statusNotification(c)
					}
				} else {
					o.Asset.Evses[x].Connectors[y].DP.Position = 0
					o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
				}
			} else {
				o.notAutoChargePoint(c, x, y)
			}

			o.Asset.Evses[x].Connectors[y].TPower = assets.CalculateTotalPower(c.TPower, c.Data[c.DP.Position].Power)
			o.Asset.Evses[x].Connectors[y].Energy = assets.CalculateEnergy(
				o.Asset.Evses[x].Connectors[y].TPower,
				o.Asset.Evses[x].Connectors[y].Energy,
				c.Data[c.DP.Position].Power,
				o.st,
			)
		}
	}
}

func (o *Ocpp16) meterValuesSampledData() {
	fmt.Println("meter values")

	for _, c := range o.Asset.Evses[0].Connectors {
		var mvl []types.MeterValue

		var spl []types.SampledValue

		var confL = strings.Split(*o.Conf["MeterValuesSampledData"].Value, ",")

		for i := 0; i < int(o.Asset.Phases); i++ {
			for _, conf := range confL {
				var sp types.SampledValue

				switch conf {
				case assets.EnergyActiveImportRegister:
					sp = types.SampledValue{
						Value:     strconv.FormatFloat(c.Energy, 'f', 4, 64),
						Unit:      types.UnitOfMeasureWh,
						Format:    types.ValueFormatRaw,
						Measurand: types.Measurand(assets.EnergyActiveImportRegister),
						Phase:     types.Phase(assets.Phases[i]),
					}

				case assets.Voltage:
					sp = types.SampledValue{
						Value:     strconv.FormatInt(c.Data[c.DP.Position].Voltage[i], 10),
						Unit:      types.UnitOfMeasureV,
						Format:    types.ValueFormatRaw,
						Measurand: types.Measurand(assets.Voltage),
						Phase:     types.Phase(assets.Phases[i]),
					}

				case assets.CurrentImport:
					sp = types.SampledValue{
						Value: strconv.FormatFloat(assets.CalculateCurrent(
							c.Data[c.DP.Position].Power,
							c.Data[c.DP.Position].PowerFactor,
							c.Data[c.DP.Position].Voltage[i],
							int64(o.Asset.Phases),
						), 'f', 4, 64),
						Unit:      types.UnitOfMeasureA,
						Format:    types.ValueFormatRaw,
						Measurand: types.Measurand(assets.CurrentImport),
						Phase:     types.Phase(assets.Phases[i]),
					}
				case assets.PowerActiveImport:
					sp = types.SampledValue{
						Value:     strconv.FormatFloat(float64(c.Data[c.DP.Position].Power), 'f', 4, 64),
						Unit:      types.UnitOfMeasureW,
						Format:    types.ValueFormatRaw,
						Measurand: types.Measurand(assets.PowerActiveImport),
						Phase:     types.Phase(assets.Phases[i]),
					}
				}

				spl = append(spl, sp)
			}
		}

		mvl = []types.MeterValue{
			{
				Timestamp:    types.NewDateTime(time.Now()),
				SampledValue: spl,
			},
		}

		var req = core.MeterValuesRequest{
			ConnectorId: int(c.ID),
			MeterValue:  mvl,
		}

		var lm = map[string]string{
			"protocol":  string(o.Asset.Protocol),
			"function":  "meterValuesSampledData",
			"feature":   req.GetFeatureName(),
			"simulator": o.Asset.Name,
			"type":      assets.Request,
		}

		o.logger.log(lm, req, assets.Info)

		resp, err := o.s.SendRequest(req)

		lm["type"] = assets.Response

		if err != nil {
			o.logger.log(lm, err, assets.Error)
		}

		o.logger.log(lm, resp, assets.Info)
	}
}

func (o *Ocpp16) statusNotification(c simulator.Connector) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "statusNotification",
		"feature":   "StatusNotification",
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = core.StatusNotificationRequest{
		ConnectorId: int(c.ID),
		ErrorCode:   core.ChargePointErrorCode(assets.ErrorCode[c.Data[c.DP.Position].ErrorCode]),
		Status:      core.ChargePointStatus(assets.Status[c.Data[c.DP.Position].ChargingState]),
	}

	lm["feature"] = req.GetFeatureName()

	o.logger.log(lm, req, assets.Info)

	cb := func(res ocpp.Response, err error) {
		var lm2 = map[string]string{
			"protocol":  string(o.Asset.Protocol),
			"function":  "statusNotification",
			"feature":   res.GetFeatureName(),
			"simulator": o.Asset.Name,
			"sender":    assets.CS,
			"type":      assets.Response,
		}

		o.logger.log(lm2, res, assets.Info)
	}

	err := o.s.SendRequestAsync(req, cb)

	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
	}
}

func (o *Ocpp16) stopTransaction(c simulator.Connector) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "stopTransaction",
		"feature":   "StatusNotification",
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = core.StopTransactionRequest{
		IdTag:         "asdasdasd",
		MeterStop:     int(c.Energy),
		Timestamp:     types.NewDateTime(time.Now()),
		TransactionId: o.chargeProfile.TransactionId,
	}

	lm["feature"] = req.GetFeatureName()

	o.logger.log(lm, req, assets.Info)

	cb := func(res ocpp.Response, err error) {
		var lm2 = map[string]string{
			"protocol":  string(o.Asset.Protocol),
			"function":  "stopTransaction",
			"feature":   res.GetFeatureName(),
			"simulator": o.Asset.Name,
			"sender":    assets.CS,
			"type":      assets.Response,
		}

		o.logger.log(lm2, res, assets.Info)
	}

	err := o.s.SendRequestAsync(req, cb)

	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
	}
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

/**/
func canEnable(cl []simulator.Connector) bool {
	for _, c := range cl {
		if c.Enabled {
			return false
		}
	}

	return true
}

/**/
func (o *Ocpp16) notAutoChargePoint(c simulator.Connector, x, y int) {
	if c.Enabled {
		if c.DP.Ticker < c.Data[c.DP.Position].Duration {
			o.Asset.Evses[x].Connectors[y].DP.Ticker++
		} else {
			o.Asset.Evses[x].Connectors[y].DP.Ticker = 0

			if assets.Status[c.Data[c.DP.Position].ChargingState] == "Finishing" {
				o.Asset.Evses[x].Connectors[y].DP.Position = 0
				o.Asset.Evses[x].Connectors[y].Enabled = false

				go o.statusNotification(o.Asset.Evses[x].Connectors[y])

				return
			}

			if c.DP.Position < int64(len(c.Data)-1) {
				o.Asset.Evses[x].Connectors[y].DP.Position++
			} else {
				o.Asset.Evses[x].Connectors[y].DP.Position = 0
			}

			for {
				if c.Data[c.DP.Position].ChargingState == int64(assets.Charging) {
					break
				}

				if c.DP.Position < int64(len(c.Data)-1) {
					o.Asset.Evses[x].Connectors[y].DP.Position++
				} else {
					o.Asset.Evses[x].Connectors[y].DP.Position = 0
				}
			}
		}
	} else {
		o.Asset.Evses[x].Connectors[y].DP.Position = 0
		o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
	}
}
