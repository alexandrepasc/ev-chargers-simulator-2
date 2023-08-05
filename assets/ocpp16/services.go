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
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

/*
Set the starting configurations for the asset.
*/
// TODO: need to review the way the conf default values are being set
func (o *Ocpp16) setStartUpConfigurations() {
	o.Conf = config

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "setConfigurations", "simulator": o.Asset.Name}, "Set startup configurations", assets.Info)

	var artr = o.Conf["AuthorizeRemoteTxRequests"]

	artr.Value = assets.GetStringPointer(strconv.FormatBool(o.Asset.AuthorizeRemote))

	o.Conf["AuthorizeRemoteTxRequests"] = artr

	var noc = o.Conf["NumberOfConnectors"]

	noc.Value = assets.GetStringPointer(strconv.Itoa(len(o.Asset.Evses[0].Connectors)))

	o.Conf["NumberOfConnectors"] = noc

	var rr = o.Conf["ResetRetries"]

	rr.Value = assets.GetStringPointer(strconv.FormatInt(assets.DefResetRetries, 10))

	o.Conf["ResetRetries"] = rr

	var stoesd = o.Conf["StopTransactionOnEVSideDisconnect"]

	stoesd.Value = assets.GetStringPointer(strconv.FormatBool(assets.DefStopTransactionOnEvSideDisconnect))

	o.Conf["StopTransactionOnEVSideDisconnect"] = stoesd

	var stoii = o.Conf["StopTransactionOnInvalidId"]

	stoii.Value = assets.GetStringPointer(strconv.FormatBool(assets.DefStopTransactionOnInvalidID))

	o.Conf["StopTransactionOnInvalidId"] = stoii

	var t = strconv.FormatInt(o.Timeout, 10)

	var cto = o.Conf["ConnectionTimeOut"]

	cto.Value = &t

	o.Conf["ConnectionTimeOut"] = cto

	var hb = o.Conf["HeartbeatInterval"]

	hb.Value = assets.GetStringPointer(strconv.FormatInt(assets.DefHeartbeatInterval, 10))

	o.Conf["HeartbeatInterval"] = hb

	o.Conf["MeterValueSampleInterval"] = core.ConfigurationKey{
		Key:      o.Conf["MeterValueSampleInterval"].Key,
		Readonly: o.Conf["MeterValueSampleInterval"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefMeterValueSampleInterval, 10)),
	}

	o.Conf["MeterValuesSampledData"] = core.ConfigurationKey{
		Key:      o.Conf["MeterValuesSampledData"].Key,
		Readonly: o.Conf["MeterValuesSampledData"].Readonly,
		Value:    assets.GetStringPointer(assets.DefMeterValuesSampledData),
	}

	o.Conf["ClockAlignedDataInterval"] = core.ConfigurationKey{
		Key:      o.Conf["ClockAlignedDataInterval"].Key,
		Readonly: o.Conf["ClockAlignedDataInterval"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefClockAlignedDataInterval, 10)),
	}

	o.Conf["MeterValuesAlignedData"] = core.ConfigurationKey{
		Key:      o.Conf["MeterValuesAlignedData"].Key,
		Readonly: o.Conf["MeterValuesAlignedData"].Readonly,
		Value:    assets.GetStringPointer(assets.DefMeterValuesAlignedData),
	}

	o.Conf["StopTxnAlignedData"] = core.ConfigurationKey{
		Key:      o.Conf["StopTxnAlignedData"].Key,
		Readonly: o.Conf["StopTxnAlignedData"].Readonly,
		Value:    assets.GetStringPointer(assets.DefStopTxnAlignedData),
	}

	o.Conf["StopTxnSampledData"] = core.ConfigurationKey{
		Key:      o.Conf["StopTxnSampledData"].Key,
		Readonly: o.Conf["StopTxnSampledData"].Readonly,
		Value:    assets.GetStringPointer(assets.DefStopTxnSampledData),
	}

	o.t = 0

	for x, e := range o.Asset.Evses {
		for y := range e.Connectors {
			o.Asset.Evses[x].Connectors[y].DP.Position = 0
			o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
			o.Asset.Evses[x].Connectors[y].TPower = 0
			o.Asset.Evses[x].Connectors[y].Energy = 0
			o.Asset.Evses[x].Connectors[y].Enabled = false
			o.Asset.Evses[x].Connectors[y].Availability = string(assets.Operative)
		}
	}

	o.localAuth.version = 0
	if !o.Asset.AuthList {
		o.localAuth.version = -1
	}
}

/*
Sends the boot notification to the central system, using the model to get the information.
*/
func (o *Ocpp16) sendBootNotification() {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "sendBootNotification",
		"feature":   core.BootNotificationFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = core.BootNotificationRequest{
		ChargePointSerialNumber: o.Mod.Ocpp.SerialNumb,
		ChargePointModel:        o.Mod.Ocpp.Model,
		ChargePointVendor:       o.Mod.Ocpp.Vendor,
		FirmwareVersion:         o.Mod.Ocpp.FwVersion,
		MeterSerialNumber:       o.Mod.Ocpp.MeterSerialNumb,
		Iccid:                   o.Mod.Ocpp.Modem.Iccid,
		Imsi:                    o.Mod.Ocpp.Modem.Imsi,
	}

	o.logger.log(lm, req, assets.Info)

	var res, err = o.s.SendRequest(req)

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
	}

	o.logger.log(lm, res.(*core.BootNotificationConfirmation), assets.Info)

	o.st = time.Now()
}

/*
Process the start transaction request and proceed with the correct logic in each case of the
authorize config.

Validates if the connector selected is not active, activate it. If no connector were selected
validate if any of them can be activated, if so activate it.
*/
// TODO: need to check the start transaction to go throw the preparing state instead of directly to charging
// TODO: need to check the best way to handle the id tag being used in the current transaction
func (o *Ocpp16) processRemoteStartTransaction(r *core.RemoteStartTransactionRequest) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "processRemoteStartTransaction",
		"feature":   "StartTransaction",
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var auth, err = strconv.ParseBool(*o.Conf["AuthorizeRemoteTxRequests"].Value)

	if err != nil {
		var lm2 = map[string]string{
			"protocol":  string(o.Asset.Protocol),
			"function":  "processRemoteStartTransaction",
			"simulator": o.Asset.Name,
		}

		o.logger.log(lm2, err, assets.Fatal)
	}

	if !auth {
		// TODO: the store of the charging profile should not be set at this point, since the validations if the session can be started are not done yet
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

					// check if the connector is available
					if c.Availability != string(assets.Operative) {
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

					o.statusNotification(&o.Asset.Evses[0].Connectors[i])

					return
				}
			} else {
				return
			}
		}

		if canEnable(o.Asset.Evses[0].Connectors) {
			for i, c := range o.Asset.Evses[0].Connectors {
				if c.Availability != string(assets.Operative) {
					continue
				}

				o.Asset.Evses[0].Connectors[i].Enabled = true

				o.Asset.Evses[0].Connectors[i].DP.Position = 2

				var req = core.StartTransactionRequest{
					ConnectorId: int(c.ID),
					IdTag:       r.IdTag,
					Timestamp:   types.NewDateTime(time.Now()),
				}

				go o.s.SendRequest(req) //nolint:errcheck // because at the moment can not handle the error since it is in a routine

				return
			}
		} else {
			return
		}
	}

	/*
		In case the AuthorizeRemoteTxRequests (AuthorizeRemote) is true, the CP will behave as it
		starts the session by it self. This behaviour is to try to authorize the tag in the local auth
		list and/or requesting the authorize request ot the CS.
	*/

	var c, i = o.getConnectorAndIndex(r.ConnectorId)

	if c.Availability != string(assets.Operative) {
		return
	}

	o.Asset.Evses[0].Connectors[i].DP.Position = 1

	o.statusNotification(&o.Asset.Evses[0].Connectors[i])

	var na = true

	if o.localAuth.version > 0 {
		for _, a := range o.localAuth.list {
			if a.IdTag == r.IdTag && a.IdTagInfo.Status == types.AuthorizationStatusAccepted {
				na = false
			}
		}
	}

	if na {
		if !o.authorize(r.IdTag) {
			o.Asset.Evses[0].Connectors[i].DP.Position = 0

			o.statusNotification(&o.Asset.Evses[0].Connectors[i])

			return
		}
	}

	o.Asset.Evses[0].Connectors[i].Enabled = true

	o.chargeProfile = r.ChargingProfile

	o.Asset.Evses[0].CIDTag = r.IdTag

	var resp = o.startTransaction(o.Asset.Evses[0].CIDTag, c)

	var stoii, errB = strconv.ParseBool(*o.Conf["StopTransactionOnInvalidId"].Value)

	if errB != nil {
		o.logger.log(lm, errB, assets.Error)
		return
	}

	if !stoii {
		return
	}

	if resp.IdTagInfo.Status != types.AuthorizationStatusAccepted {
		o.Asset.Evses[0].Connectors[i].DP.Position = 0
		o.Asset.Evses[0].Connectors[i].Enabled = false
		o.txnAlignedData = []types.MeterValue{}
		o.txnSampledData = []types.MeterValue{}
		o.statusNotification(&o.Asset.Evses[0].Connectors[i])
	}
}

/*
Handles the logic for the remote stop transaction filtering the transaction ID, and the active
connector to stop the session.

In case a session is stopped sends a status notification with the new connector status, the
stop transaction request, and the accepted response.

If none of the filters pass the response will be rejected.
*/
func (o *Ocpp16) processRemoteStopTransaction(r *core.RemoteStopTransactionRequest) *core.RemoteStopTransactionConfirmation {
	if r.TransactionId == o.chargeProfile.TransactionId {
		// at the moment this is only supporting 1 evse per simulator, so this will only look for one position
		for i, c := range o.Asset.Evses[0].Connectors {
			if !c.Enabled {
				continue
			}

			o.Asset.Evses[0].Connectors[i].DP.Position = int64(len(c.Data) - 1)
			o.Asset.Evses[0].Connectors[i].DP.Ticker = 0

			go o.statusNotification(&o.Asset.Evses[0].Connectors[i])

			go o.stopTransaction(o.Asset.Evses[0].CIDTag, &o.Asset.Evses[0].Connectors[i])

			return &core.RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}
		}
	}

	return &core.RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}
}

/*
Has the reset logic for the soft and hard reset.

If is a soft and a connector is enabled will set the data position to 0, disable it and send a stop
transaction request.

In a hard reset will set all the connectors data and the asset data.
*/
func (o *Ocpp16) processReset(r *core.ResetRequest) *core.ResetConfirmation {
	if r.Type == core.ResetType(assets.Soft) {
		for x, e := range o.Asset.Evses {
			for y, c := range e.Connectors {
				if !c.Enabled {
					continue
				}

				o.Asset.Evses[x].Connectors[y].Enabled = false
				o.Asset.Evses[x].Connectors[y].DP.Position = 0
				o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
				o.txnAlignedData = []types.MeterValue{}
				o.txnSampledData = []types.MeterValue{}

				go o.stopTransaction(o.Asset.Evses[x].CIDTag, &o.Asset.Evses[x].Connectors[y])
			}
		}

		return &core.ResetConfirmation{Status: core.ResetStatusAccepted}
	}

	for x, e := range o.Asset.Evses {
		for y := range e.Connectors {
			o.Asset.Evses[x].Connectors[y].Enabled = false
			o.Asset.Evses[x].Connectors[y].DP.Position = 0
			o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
			o.Asset.Evses[x].Connectors[y].TPower = 0
			o.Asset.Evses[x].Connectors[y].Energy = 0
			o.Asset.Evses[x].Connectors[y].Availability = string(assets.Operative)
		}
	}

	o.localAuth.version = 0
	o.localAuth.list = nil
	o.chargeProfile = nil
	o.t = 0

	go o.sendBootNotification()

	return &core.ResetConfirmation{Status: core.ResetStatusAccepted}
}

/*
Filter the connector in the request, change it's availability, and sends an accepted response.

In case the connector has a session it will change the availability, but the response will be
scheduled.

If the connector sent doesn't match the CP connectors will return rejected.
*/
// TODO: does not support the 0 connector logic that would change the state of the CP and all it's connectors
func (o *Ocpp16) processChangeAvailability(r *core.ChangeAvailabilityRequest) *core.ChangeAvailabilityConfirmation {
	// at the moment it is only supporting one evse per simulator so the evse will be set to 0 position
	for i, c := range o.Asset.Evses[0].Connectors {
		if r.ConnectorId != int(c.ID) {
			continue
		}

		var aux simulator.Connector

		// if the availability change matches the curren it returns accept directly
		if o.Asset.Evses[0].Connectors[i].Availability == string(r.Type) {
			return &core.ChangeAvailabilityConfirmation{Status: core.AvailabilityStatusAccepted}
		}

		o.Asset.Evses[0].Connectors[i].Availability = string(r.Type)

		// if the request changes the availability and the connector is charging returns scheduled
		if c.Data[c.DP.Position].ChargingState == int64(assets.Charging) {
			return &core.ChangeAvailabilityConfirmation{Status: core.AvailabilityStatusScheduled}
		}

		aux = o.Asset.Evses[0].Connectors[i]

		switch o.Asset.Evses[0].Connectors[i].Availability {
		case string(assets.Inoperative):
			aux.Data[aux.DP.Position].ChargingState = int64(assets.Unavailable)

		case string(assets.Operative):
			aux.Data[aux.DP.Position].ChargingState = int64(assets.Available)
		}

		go o.statusNotification(&aux)

		return &core.ChangeAvailabilityConfirmation{Status: core.AvailabilityStatusAccepted}
	}

	return &core.ChangeAvailabilityConfirmation{Status: core.AvailabilityStatusRejected}
}

/*
Logic to handle the update of the local authorization request from the CS. Checks if the update
type is full or differential and behaves according.

Validates if there is no problem with the update and fails the request in case of a failure.
*/
func (o *Ocpp16) processSendLocalList(r *localauth.SendLocalListRequest) *localauth.SendLocalListConfirmation {
	if r.UpdateType == localauth.UpdateTypeFull {
		o.localAuth.version = int64(r.ListVersion)
		o.localAuth.list = r.LocalAuthorizationList

		return &localauth.SendLocalListConfirmation{Status: localauth.UpdateStatusAccepted}
	}

	if r.ListVersion <= int(o.localAuth.version) {
		return &localauth.SendLocalListConfirmation{Status: localauth.UpdateStatusVersionMismatch}
	}

	var al []localauth.AuthorizationData

	for i := range r.LocalAuthorizationList {
		if i < len(o.localAuth.list) {
			if r.LocalAuthorizationList[i].IdTag == o.localAuth.list[i].IdTag {
				al = append(al, r.LocalAuthorizationList[i])
			} else {
				return &localauth.SendLocalListConfirmation{Status: localauth.UpdateStatusFailed}
			}
		} else {
			al = append(al, r.LocalAuthorizationList[i])
		}
	}

	o.localAuth.list = al

	return &localauth.SendLocalListConfirmation{Status: localauth.UpdateStatusAccepted}
}

/*
Process the logic to trigger the sampled data meter values. Get the configured interval, check
if it is 0, check if it is time to send the message, gets the configuration data, and evaluate
if any of the connectors is active to send the request. If a connector is active get the
information and call the send function.
*/
func (o *Ocpp16) processSampledData() {
	var v, errI = strconv.ParseInt(*o.Conf["MeterValueSampleInterval"].Value, 10, 64)

	if errI != nil {
		var lm2 = map[string]string{
			"protocol":  string(o.Asset.Protocol),
			"function":  "processSampledData",
			"simulator": o.Asset.Name,
		}

		o.logger.log(lm2, errI, assets.Fatal)

		return
	}

	if v == 0 {
		return
	}

	if o.t%v != 0 {
		return
	}

	var conf = strings.Split(*o.Conf["MeterValuesSampledData"].Value, ",")

	for ie, e := range o.Asset.Evses {
		for ic, c := range e.Connectors {
			if !c.Enabled {
				continue
			}

			if c.Data[c.DP.Position].ChargingState != int64(assets.Charging) {
				continue
			}

			var sd = o.meterValuesSampledData(ie, ic, conf)

			var cl = strings.Split(*o.Conf["StopTxnSampledData"].Value, ",")

			if len(cl) != 0 || cl[0] != "" {
				var tsd = o.meterValuesSampledData(ie, ic, cl)

				o.txnSampledData = append(o.txnSampledData, tsd...)
			}

			o.meterValues(c.ID, sd)
		}
	}
}

/*
Handles the logic to send the sampled data message to the CS when triggered by the
trigger message request. Will loop by the connectors if a connector id match build the meter
values message and send it to the CS.

id	-	Connector identifier sent by the CS request
*/
func (o *Ocpp16) processTriggerSampledData(id *int) {
	// There is no way to identify the evse so ir will be set as 0 the array index
	for ic, c := range o.Asset.Evses[0].Connectors {
		if c.ID != int64(*id) {
			continue
		}

		var conf = strings.Split(*o.Conf["MeterValuesSampledData"].Value, ",")

		var sd = o.meterValuesSampledData(0, ic, conf)

		o.meterValues(c.ID, sd)

		break
	}
}

/*
Builds the meter value sampled data list with the configuration values for the connector. It will
return the list with the data for each phase ([]types.MeterValue). The supported information keys
that can be used in this request are specified in the constants file.

ie		-	Evse array index (int)

ic		-	Connector array index (int)

confL	-	Configuration list with the data needed to the request ([]string)
*/
func (o *Ocpp16) meterValuesSampledData(ie, ic int, confL []string) []types.MeterValue {
	var spl []types.SampledValue

	for _, conf := range confL {
		var confT = strings.TrimSpace(conf)

		for i := 0; i < int(o.Asset.Phases); i++ {
			var sp types.SampledValue

			var cdp = o.Asset.Evses[ie].Connectors[ic].DP.Position

			switch confT {
			case assets.EnergyActiveImportRegister:
				sp = types.SampledValue{
					Value:     strconv.FormatFloat(o.Asset.Evses[ie].Connectors[ic].Energy, 'f', 4, 64),
					Unit:      types.UnitOfMeasureWh,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.EnergyActiveImportRegister),
					Phase:     types.Phase(assets.Phases[i]),
				}

			// TODO: this is not being calculated and the value is set to 0
			case assets.EnergyReactiveImportRegister:
				sp = types.SampledValue{
					Value:     "0",
					Unit:      types.UnitOfMeasureVarh,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.EnergyReactiveImportRegister),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.Voltage:
				sp = types.SampledValue{
					Value:     strconv.FormatInt(o.Asset.Evses[ie].Connectors[ic].Data[cdp].Voltage[i], 10),
					Unit:      types.UnitOfMeasureV,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.Voltage),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.CurrentImport:
				sp = types.SampledValue{
					Value: strconv.FormatFloat(assets.CalculateCurrent(
						o.Asset.Evses[ie].Connectors[ic].Data[cdp].Power,
						o.Asset.Evses[ie].Connectors[ic].Data[cdp].PowerFactor,
						o.Asset.Evses[ie].Connectors[ic].Data[cdp].Voltage[i],
						int64(o.Asset.Phases),
					), 'f', 4, 64),
					Unit:      types.UnitOfMeasureA,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.CurrentImport),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.PowerActiveImport:
				sp = types.SampledValue{
					Value:     strconv.FormatFloat(float64(o.Asset.Evses[ie].Connectors[ic].Data[cdp].Power), 'f', 4, 64),
					Unit:      types.UnitOfMeasureW,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.PowerActiveImport),
					Phase:     types.Phase(assets.Phases[i]),
				}
			}

			spl = append(spl, sp)
		}
	}

	var mvl = []types.MeterValue{
		{
			Timestamp:    types.NewDateTime(time.Now()),
			SampledValue: spl,
		},
	}

	return mvl
}

/*
Process the logic to trigger the sampled data meter values. Get the configured interval, check
if it is 0, check if it is time to send the message, gets the configuration data, and call the send
function.

If a connector is active get the transaction aligned data values and append them in a variable, to
be used in the stop transaction message.
*/
func (o *Ocpp16) processAlignedData() {
	var i = time.Now().UTC().Sub(time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC)).Seconds()

	var v, errV = strconv.ParseFloat(*o.Conf["ClockAlignedDataInterval"].Value, 32)

	if errV != nil {
		var lm2 = map[string]string{
			"protocol":  string(o.Asset.Protocol),
			"function":  "processAlignedData",
			"simulator": o.Asset.Name,
		}

		o.logger.log(lm2, errV, assets.Fatal)

		return
	}

	if v == 0 {
		return
	}

	if int64(i)%int64(v) != 0 {
		return
	}

	var conf = strings.Split(*o.Conf["MeterValuesAlignedData"].Value, ",")

	var ad = o.meterValuesAlignedData(conf)

	for _, e := range o.Asset.Evses {
		if canEnable(e.Connectors) {
			continue
		}

		var _, c = getActiveConnector(e)

		if c.Data[c.DP.Position].ChargingState != int64(assets.Charging) {
			continue
		}

		var cl = strings.Split(*o.Conf["StopTxnAlignedData"].Value, ",")

		if len(cl) == 0 || cl[0] == "" {
			continue
		}

		var tad = o.meterValuesAlignedData(cl)

		o.txnAlignedData = append(o.txnAlignedData, tad...)
	}

	// TODO: review the connector id, at this moment is returning the total of the asset so the id is 0
	o.meterValues(0, ad)
}

/*
Builds the meter value aligned data list with the configuration values for the asset. It will
return the list with the data ([]types.MeterValue). The supported information keys that can be used
in this request are specified in the constants file.

confL	-	Configuration list with the data needed to the request ([]string)
*/
// TODO: some more research is needed to this functionality
func (o *Ocpp16) meterValuesAlignedData(confL []string) []types.MeterValue {
	var spl = []types.SampledValue{}

	var tp = assets.CalculateCPPower(o.Asset.Evses)

	for _, conf := range confL {
		var confT = strings.TrimSpace(conf)

		var sp types.SampledValue

		switch confT {
		case assets.EnergyActiveImportRegister:
			sp = types.SampledValue{
				Value:     strconv.FormatFloat(assets.CalculateCPEnergy(tp, o.st), 'f', 4, 64),
				Unit:      types.UnitOfMeasureWh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyActiveImportRegister),
			}

		case assets.EnergyReactiveImportRegister:
			sp = types.SampledValue{
				Value:     "0",
				Unit:      types.UnitOfMeasureVarh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyReactiveImportRegister),
			}

		case assets.PowerActiveImport:
			sp = types.SampledValue{
				Value:     strconv.FormatFloat(tp, 'f', 4, 64),
				Unit:      types.UnitOfMeasureW,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.PowerActiveImport),
			}
		}

		spl = append(spl, sp)
	}

	var mvl = []types.MeterValue{
		{
			Timestamp:    types.NewDateTime(time.Now()),
			SampledValue: spl,
		},
	}

	return mvl
}

/*
Updates the connectors data position and ticker, it filters if the simulator has the start charging
at true (automatic) or false (passive) to do these logic.

In automatic mode if no connector is enabled will enable one, if the connector is enabled will
do the logic of passing on all the data that is defined in the connector data configuration file.
In case the connector data changes state it will send the status notification to the CS. If the
data charging state is Finishing (6) sends the stop transaction request. If the connector is not
enable sets the data position and ticker to 0.

If the start charging in the simulator is false (passive) it will call the notAutoChargePoint
function to handle the logic.
*/
// TODO: the total power calculation need to be reviewed, at the moment with 100 w in a couple of secs the result is 0
func (o *Ocpp16) updateData() {
	for x, e := range o.Asset.Evses {
		for y, c := range e.Connectors {
			if o.Asset.StartCharging {
				if canEnable(e.Connectors) {
					o.Asset.Evses[x].Connectors[y].Enabled = true

					// hard code the transaction id to the start charging simulator
					o.chargeProfile = &types.ChargingProfile{
						TransactionId: 1,
					}
				}

				if o.Asset.Evses[x].Connectors[y].Enabled {
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

					if cs != c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState {
						o.statusNotification(&o.Asset.Evses[x].Connectors[y])

						if c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState == int64(assets.Finishing) {
							o.stopTransaction(o.Asset.Evses[x].CIDTag, &o.Asset.Evses[x].Connectors[y])
							o.txnAlignedData = []types.MeterValue{}
							o.txnSampledData = []types.MeterValue{}
						}

						// If the connector starts charging send the start transaction request
						if c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState == int64(assets.Charging) {
							// TODO: need to review the id tag
							var resp = o.startTransaction("QWEASDZXC", &o.Asset.Evses[x].Connectors[y])

							var lm = map[string]string{
								"protocol":  string(o.Asset.Protocol),
								"function":  "updateData",
								"feature":   resp.GetFeatureName(),
								"simulator": o.Asset.Name,
								"sender":    assets.CS,
								"type":      assets.Response,
							}

							var stoii, errB = strconv.ParseBool(*o.Conf["StopTransactionOnInvalidId"].Value)

							if errB != nil {
								o.logger.log(lm, errB, assets.Error)
								return
							}

							if stoii {
								if resp.IdTagInfo.Status != types.AuthorizationStatusAccepted {
									o.logger.log(lm, resp, assets.Info)

									for i := c.DP.Position; i < int64(len(c.Data)); i++ {
										if c.Data[i].ChargingState == int64(assets.Finishing) {
											o.Asset.Evses[x].Connectors[y].DP.Position = i
											o.Asset.Evses[x].Connectors[y].DP.Ticker = 0

											break
										}
									}
								}
							}
						}
					}
				} else {
					o.Asset.Evses[x].Connectors[y].DP.Position = 0
					o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
				}
			} else {
				o.notAutoChargePoint(&o.Asset.Evses[x].Connectors[y], x, y)
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

/*
Execute the status notification request with the current charging state of the connector.

c	-	Connector structure with all it's data (*simulator.Connector)
*/
func (o *Ocpp16) statusNotification(c *simulator.Connector) {
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

/*
id	-	The tag id used to authorize the session (string)

c	-	The connector that will be used in the session (*simulator.Connector)
*/
func (o *Ocpp16) startTransaction(id string, c *simulator.Connector) *core.StartTransactionConfirmation {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "startTransaction",
		"feature":   core.StartTransactionFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = core.StartTransactionRequest{
		ConnectorId: int(c.ID),
		IdTag:       id,
		MeterStart:  int(c.Energy),
		Timestamp:   types.NewDateTime(time.Now()),
	}

	o.logger.log(lm, req, assets.Info)

	var res, err = o.s.SendRequest(req)

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
	}

	o.logger.log(lm, res.(*core.StartTransactionConfirmation), assets.Info)

	return res.(*core.StartTransactionConfirmation)
}

/*
Sends the stop transaction request for the connector.

id	-	Session id tag (string)

c	-	Evse connector information (*simulator.Connector)
*/
func (o *Ocpp16) stopTransaction(id string, c *simulator.Connector) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "stopTransaction",
		"feature":   "StatusNotification",
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var td = o.txnAlignedData
	td = append(td, o.txnSampledData...)

	// TODO: the id tag needs to be created and stored and not hard coded
	var req = core.StopTransactionRequest{
		IdTag:           id,
		MeterStop:       int(c.Energy),
		Timestamp:       types.NewDateTime(time.Now()),
		TransactionId:   o.chargeProfile.TransactionId,
		TransactionData: td,
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

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
	}
}

/*
Sends the authorize request to the CS and returns true (bool) if the id was accepted, and false if
not.

id	-	The user tag id that tries to start the session (string)
*/
func (o *Ocpp16) authorize(id string) bool {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "authorize",
		"feature":   core.AuthorizeFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = core.AuthorizeRequest{
		IdTag: id,
	}

	o.logger.log(lm, req, assets.Info)

	res, err := o.s.SendRequest(req)

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
		return false
	}

	o.logger.log(lm, res.(*core.AuthorizeConfirmation), assets.Info)

	return res.(*core.AuthorizeConfirmation).IdTagInfo.Status == types.AuthorizationStatusAccepted
}

/*
Send the heartbeat request to the CS.
*/
func (o *Ocpp16) heartbeat() {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "heartbeat",
		"feature":   core.HeartbeatFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = core.HeartbeatRequest{}

	o.logger.log(lm, req, assets.Info)

	res, err := o.s.SendRequest(req)

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
	}

	o.logger.log(lm, res.(*core.HeartbeatConfirmation), assets.Info)
}

/*
id	-	Connector identifier (int64)

mvl	-	List of meter values to sent ([]types.MeterValue)
*/
func (o *Ocpp16) meterValues(id int64, mvl []types.MeterValue) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "meterValues",
		"feature":   "",
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = core.MeterValuesRequest{
		ConnectorId: int(id),
		MeterValue:  mvl,
	}

	lm["feature"] = req.GetFeatureName()

	o.logger.log(lm, req, assets.Info)

	resp, err := o.s.SendRequest(req)

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
	}

	o.logger.log(lm, resp, assets.Info)
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

/*
Handles the counter from the simulator and handles the max int64 value, in case it is reaching the
max value (max int64 value - 7) it will be reseted to 0.
*/
func (o *Ocpp16) handleTick() {
	const rInt64 = math.MaxInt64 - 7

	if o.t > rInt64 {
		o.t = 0
	} else {
		o.t++
	}
}

/*
Checks all the connectors from an EVSE and returs false if any of the connectors are enabled.
If no connector is enabled it returns true (bool).

cl	-	EVSE connectors list ([]simulator.Connector)
*/
func canEnable(cl []simulator.Connector) bool {
	for _, c := range cl {
		if c.Enabled {
			return false
		}
	}

	return true
}

/*
Handles the update data for a simulator that has the start charging with false. It will only
update the data position and ticker in case the connector is enabled. In case the charging state
changes send the status notification request. If the charging state is Finishing (6) it will
evaluate the availability and send the status notification with the corresponding charging state.

If the connector is not enabled resets the data position and ticker to 0.

c	-	Evse connector information (*simulator.Connector)

x	-	Evse index position (int)

y	-	Evse connector index position (int)
*/
func (o *Ocpp16) notAutoChargePoint(c *simulator.Connector, x, y int) {
	if c.Enabled {
		if c.DP.Ticker < c.Data[c.DP.Position].Duration {
			o.Asset.Evses[x].Connectors[y].DP.Ticker++
		} else {
			o.Asset.Evses[x].Connectors[y].DP.Ticker = 0

			var cs = c.Data[c.DP.Position].ChargingState

			if assets.Status[c.Data[c.DP.Position].ChargingState] == assets.Status[assets.Finishing] {
				o.Asset.Evses[x].Connectors[y].DP.Position = 0
				o.Asset.Evses[x].Connectors[y].Enabled = false
				o.txnAlignedData = []types.MeterValue{}
				o.txnSampledData = []types.MeterValue{}

				var aux = o.Asset.Evses[x].Connectors[y]

				switch o.Asset.Evses[x].Connectors[y].Availability {
				case string(assets.Inoperative):
					aux.Data[aux.DP.Position].ChargingState = int64(assets.Unavailable)

				case string(assets.Operative):
					aux.Data[aux.DP.Position].ChargingState = int64(assets.Available)
				}

				go o.statusNotification(&aux)

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

			if cs != c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState {
				o.statusNotification(&o.Asset.Evses[x].Connectors[y])

				if c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState == int64(assets.Finishing) {
					o.stopTransaction(o.Asset.Evses[x].CIDTag, &o.Asset.Evses[x].Connectors[y])
				}
			}
		}
	} else {
		o.Asset.Evses[x].Connectors[y].DP.Position = 0
		o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
		o.Asset.Evses[x].Connectors[y].Enabled = false
		o.txnAlignedData = []types.MeterValue{}
		o.txnSampledData = []types.MeterValue{}
	}
}

/*
Get the connector object that match the id argument value. It will return the object and the
connector index from the asset object.

If the id value doesn't match any of the evse connectors it will return a nil objet and the index
with the 0 value.

In case the id argument has the value nil it will return the first connector that is not set as
enabled.

id	-	The id returned by the CS request (*int)
*/
func (o *Ocpp16) getConnectorAndIndex(id *int) (ci *simulator.Connector, index int) {
	if id != nil {
		for i, c := range o.Asset.Evses[0].Connectors {
			if c.ID == int64(*id) {
				return &c, i
			}
		}

		return nil, 0
	}

	for i, c := range o.Asset.Evses[0].Connectors {
		if !c.Enabled {
			return &c, i
		}
	}

	return nil, 0
}

/*
Get from the evse the connecto that is active, in case one of the connectors is active returns
the evse connector index (int) and the connector structure (*simulator.Connector).

In case none of the connectors is active will return the index -1 and the structure as nil.

e	-	Evse structure simulator.Evse
*/
func getActiveConnector(e simulator.Evse) (i int, c *simulator.Connector) {
	for ci, c := range e.Connectors {
		if c.Enabled {
			return ci, &c
		}
	}

	return -1, nil
}
