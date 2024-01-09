package ocpp16

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

/*
Set the starting configurations for the asset.
*/
func (o *Ocpp16) setStartUpConfigurations() {
	o.Conf = config

	o.connectSeq = true

	o.bootSeq = assets.BootSeq{
		IsToTrigger:  true,
		BootStatus:   core.RegistrationStatusAccepted,
		BootReason:   nil,
		BootInterval: 0,
	}

	o.disconnectSeq = false

	o.resetSeq = assets.ResetSeq{
		IsToTrigger: false,
		EvseIndex:   nil,
	}

	o.heartbeatC = 0

	o.logger.log(map[string]string{"protocol": string(o.Asset.Protocol), "function": "setStartUpConfigurations", "simulator": o.Asset.Name},
		o.L.Get(text.StartUpConfigurations), assets.Info)

	// This will load only the password since there is no information in the documentation regarding how to set the user in the cp
	// TODO: do not think that this is correct another review to this should be made
	// if o.Asset.BasicAuth {
	// 	o.Conf["AuthorizationKey"] = core.ConfigurationKey{
	// 		Key:      o.Conf["AuthorizationKey"].Key,
	// 		Readonly: o.Conf["AuthorizationKey"].Readonly,
	// 		Value:    assets.GetStringPointer(o.Mod.BasicAuth.Password),
	// 	}
	// }

	o.Conf["AuthorizeRemoteTxRequests"] = core.ConfigurationKey{
		Key:      o.Conf["AuthorizeRemoteTxRequests"].Key,
		Readonly: o.Conf["AuthorizeRemoteTxRequests"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatBool(o.Asset.AuthorizeRemote)),
	}

	o.Conf["NumberOfConnectors"] = core.ConfigurationKey{
		Key:      o.Conf["NumberOfConnectors"].Key,
		Readonly: o.Conf["NumberOfConnectors"].Readonly,
		Value:    assets.GetStringPointer(strconv.Itoa(len(o.Asset.Evses[0].Connectors))),
	}

	o.Conf["ResetRetries"] = core.ConfigurationKey{
		Key:      o.Conf["ResetRetries"].Key,
		Readonly: o.Conf["ResetRetries"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefResetRetries, 10)),
	}

	o.Conf["StopTransactionOnEVSideDisconnect"] = core.ConfigurationKey{
		Key:      o.Conf["StopTransactionOnEVSideDisconnect"].Key,
		Readonly: o.Conf["StopTransactionOnEVSideDisconnect"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatBool(assets.DefStopTransactionOnEvSideDisconnect)),
	}

	o.Conf["StopTransactionOnInvalidId"] = core.ConfigurationKey{
		Key:      o.Conf["StopTransactionOnInvalidId"].Key,
		Readonly: o.Conf["StopTransactionOnInvalidId"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatBool(assets.DefStopTransactionOnInvalidID)),
	}

	o.Conf["ConnectionTimeOut"] = core.ConfigurationKey{
		Key:      o.Conf["ConnectionTimeOut"].Key,
		Readonly: o.Conf["ConnectionTimeOut"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(o.Timeout, 10)),
	}

	o.Conf["HeartbeatInterval"] = core.ConfigurationKey{
		Key:      o.Conf["HeartbeatInterval"].Key,
		Readonly: o.Conf["HeartbeatInterval"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefHeartbeatInterval, 10)),
	}

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

	o.Conf["TransactionMessageAttempts"] = core.ConfigurationKey{
		Key:      o.Conf["TransactionMessageAttempts"].Key,
		Readonly: o.Conf["TransactionMessageAttempts"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefTransactionMessageAttempts, 10)),
	}

	o.Conf["TransactionMessageRetryInterval"] = core.ConfigurationKey{
		Key:      o.Conf["TransactionMessageRetryInterval"].Key,
		Readonly: o.Conf["TransactionMessageRetryInterval"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefTransactionMessageRetryInterval, 10)),
	}

	o.Conf["UnlockConnectorOnEVSideDisconnect"] = core.ConfigurationKey{
		Key:      o.Conf["UnlockConnectorOnEVSideDisconnect"].Key,
		Readonly: o.Conf["UnlockConnectorOnEVSideDisconnect"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatBool(assets.DefUnlockConnectorOnEVSideDisconnect)),
	}

	o.Conf["ConnectorPhaseRotation"] = core.ConfigurationKey{
		Key:      o.Conf["ConnectorPhaseRotation"].Key,
		Readonly: o.Conf["ConnectorPhaseRotation"].Readonly,
		Value:    (*string)(&o.Asset.PhaseRotation),
	}

	o.Conf["GetConfigurationMaxKeys"] = core.ConfigurationKey{
		Key:      o.Conf["GetConfigurationMaxKeys"].Key,
		Readonly: o.Conf["GetConfigurationMaxKeys"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefGetConfigurationMaxKeys, 10)),
	}

	o.Conf["LocalAuthorizeOffline"] = core.ConfigurationKey{
		Key:      o.Conf["LocalAuthorizeOffline"].Key,
		Readonly: o.Conf["LocalAuthorizeOffline"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatBool(assets.DefLocalAuthorizeOffline)),
	}

	o.Conf["LocalPreAuthorize"] = core.ConfigurationKey{
		Key:      o.Conf["LocalPreAuthorize"].Key,
		Readonly: o.Conf["LocalPreAuthorize"].Readonly,
		Value:    assets.GetStringPointer(strconv.FormatBool(assets.DefLocalPreAuthorize)),
	}

	o.tick = 0

	for y := range o.Asset.Evses[0].Connectors {
		o.Asset.Evses[0].Connectors[y].DP.Position = 0
		o.Asset.Evses[0].Connectors[y].DP.Ticker = 0
		o.Asset.Evses[0].Connectors[y].TPower = 0
		o.Asset.Evses[0].Connectors[y].Energy = 0
		o.Asset.Evses[0].Connectors[y].Enabled = false
		o.Asset.Evses[0].Connectors[y].Availability = string(assets.Operative)
		o.Asset.Evses[0].Connectors[y].CurrentSoC = 0
	}

	o.localAuth.version = 0
	if !o.Asset.AuthList {
		o.localAuth.version = -1
	}
}

/*
Sends the boot notification to the central system, using the model to get the information.
*/
func (o *Ocpp16) sendBootNotification() (resp *core.BootNotificationConfirmation, err error) {
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

	var res, e = o.s.SendRequest(req)

	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if e != nil {
		o.logger.log(lm, e, assets.Error)

		return nil, e
	}

	o.logger.log(lm, res.(*core.BootNotificationConfirmation), assets.Info)

	o.bootSeq.BootStatus = res.(*core.BootNotificationConfirmation).Status

	return res.(*core.BootNotificationConfirmation), nil
}

/*
Have the logic needed to process the boot notification response.

res	-	Boot notification response from the cs (*core.BootNotificationConfirmation)
*/
func (o *Ocpp16) processBootResponse(res *core.BootNotificationConfirmation) {
	if res.Status != core.RegistrationStatusAccepted {
		o.bootSeq.BootInterval = res.Interval

		if res.Interval <= 0 {
			o.bootSeq.BootInterval = int(assets.DefHeartbeatInterval)
		}
	} else {
		if res.Interval > 0 {
			o.Conf["HeartbeatInterval"] = core.ConfigurationKey{
				Key:      o.Conf["HeartbeatInterval"].Key,
				Readonly: o.Conf["HeartbeatInterval"].Readonly,
				Value:    assets.GetStringPointer(strconv.FormatInt(int64(res.Interval), 10)),
			}
		} else {
			o.Conf["HeartbeatInterval"] = core.ConfigurationKey{
				Key:      o.Conf["HeartbeatInterval"].Key,
				Readonly: o.Conf["HeartbeatInterval"].Readonly,
				Value:    assets.GetStringPointer(strconv.FormatInt(assets.DefHeartbeatInterval, 10)),
			}
		}

		// TODO: Internal clock synchronization needs to be done

		o.bootSeq.BootInterval = 0

		o.bootSeq.IsToTrigger = false
	}
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

					o.startTransaction(r.IdTag, &o.Asset.Evses[0].Connectors[i])

					o.sendStatusNotification(&o.Asset.Evses[0].Connectors[i])

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

				o.startTransaction(strconv.FormatInt(c.ID, 10), &o.Asset.Evses[0].Connectors[i])

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

	o.sendStatusNotification(&o.Asset.Evses[0].Connectors[i])

	var na = true

	if assets.GetBoolFromString(*o.Conf["LocalPreAuthorize"].Value) {
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

				o.sendStatusNotification(&o.Asset.Evses[0].Connectors[i])

				return
			}
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
		o.sendStatusNotification(&o.Asset.Evses[0].Connectors[i])
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

			go o.sendStatusNotification(&o.Asset.Evses[0].Connectors[i])

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
func (o *Ocpp16) processResetRequest(r *core.ResetRequest) *core.ResetConfirmation {
	if r.Type == core.ResetType(assets.Soft) {
		// if there is no active connectors
		if canEnable(o.Asset.Evses[0].Connectors) {
			o.disconnectSeq = true

			o.connectSeq = true

			o.bootSeq.BootInterval = 0
			o.bootSeq.IsToTrigger = true

			return &core.ResetConfirmation{Status: core.ResetStatusAccepted}
		}

		for y, c := range o.Asset.Evses[0].Connectors {
			if !c.Enabled {
				continue
			}

			const dpr = 2

			// set the duration to the end of the last state before finish so the update data will trigger the stop transaction call
			o.Asset.Evses[0].Connectors[y].DP.Position = int64(len(o.Asset.Evses[0].Connectors[y].Data) - dpr)
			o.Asset.Evses[0].Connectors[y].DP.Ticker = o.Asset.Evses[0].Connectors[y].Data[o.Asset.Evses[0].Connectors[y].DP.Position].Duration
		}

		o.resetSeq.IsToTrigger = true

		return &core.ResetConfirmation{Status: core.ResetStatusAccepted}
	}

	// if the reset has the type hard
	for y := range o.Asset.Evses[0].Connectors {
		o.Asset.Evses[0].Connectors[y].Enabled = false
		o.Asset.Evses[0].Connectors[y].DP.Position = 0
		o.Asset.Evses[0].Connectors[y].DP.Ticker = 0
		o.Asset.Evses[0].Connectors[y].CurrentSoC = 0
		o.Asset.Evses[0].Connectors[y].TPower = 0
		o.Asset.Evses[0].Connectors[y].Energy = 0
		o.Asset.Evses[0].Connectors[y].Availability = string(assets.Operative)
	}

	o.localAuth.version = 0
	o.localAuth.list = nil
	o.chargeProfile = nil

	o.disconnectSeq = true

	o.connectSeq = true

	o.bootSeq.IsToTrigger = true
	o.bootSeq.BootInterval = 0

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

		go o.sendStatusNotification(&aux)

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
	if o.bootSeq.BootStatus.(core.RegistrationStatus) != core.RegistrationStatusAccepted {
		return
	}

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

	if o.tick%v != 0 {
		return
	}

	var conf = strings.Split(*o.Conf["MeterValuesSampledData"].Value, ",")

	for ic, c := range o.Asset.Evses[0].Connectors {
		if !c.Enabled {
			continue
		}

		if c.Data[c.DP.Position].ChargingState != int64(assets.Charging) {
			continue
		}

		var sd = o.meterValuesSampledData(0, ic, conf)

		var cl = strings.Split(*o.Conf["StopTxnSampledData"].Value, ",")

		if len(cl) != 0 || cl[0] != "" {
			var tsd = o.meterValuesSampledData(0, ic, cl)

			o.txnSampledData = append(o.txnSampledData, tsd...)
		}

		o.meterValues(c.ID, sd)
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
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "meterValuesSampledData",
		"simulator": o.Asset.Name,
	}

	var spl []types.SampledValue

	for _, conf := range confL {
		var confT = strings.TrimSpace(conf)

		for i := 0; i < int(o.Asset.Phases); i++ {
			var sp types.SampledValue

			var cdp = o.Asset.Evses[ie].Connectors[ic].DP.Position

			switch confT {
			case assets.CurrentExport:
				sp = types.SampledValue{
					Value: strconv.FormatFloat(assets.CalculateCurrent(
						o.Asset.Evses[ie].Connectors[ic].Data[cdp].PowerExport,
						o.Asset.Evses[ie].Connectors[ic].Data[cdp].PowerFactor,
						o.Asset.Evses[ie].Connectors[ic].Data[cdp].Voltage[i],
						int64(o.Asset.Phases),
					), 'f', 4, 64),
					Unit:      types.UnitOfMeasureA,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.CurrentExport),
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

			case assets.CurrentOffered:
				sp = types.SampledValue{
					Value: strconv.FormatFloat(assets.CalculateConnectorMaxCurrent(
						&o.Asset.Evses[ie].Connectors[ic],
						int(o.Asset.Phases),
						i,
					), 'f', 4, 64),
					Unit:      types.UnitOfMeasureA,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.CurrentOffered),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.EnergyActiveExportRegister:
				sp = types.SampledValue{
					Value:     strconv.FormatFloat(o.Asset.Evses[ie].Connectors[ic].EnergyExport, 'f', 4, 64),
					Unit:      types.UnitOfMeasureWh,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.EnergyActiveExportRegister),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.EnergyActiveImportRegister:
				sp = types.SampledValue{
					Value:     strconv.FormatFloat(o.Asset.Evses[ie].Connectors[ic].Energy, 'f', 4, 64),
					Unit:      types.UnitOfMeasureWh,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.EnergyActiveImportRegister),
					Phase:     types.Phase(assets.Phases[i]),
				}

			// TODO: this is not being calculated and the value is set to 0
			case assets.EnergyReactiveExportRegister:
				sp = types.SampledValue{
					Value:     "0",
					Unit:      types.UnitOfMeasureVarh,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.EnergyReactiveExportRegister),
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

			case assets.EnergyActiveExportInterval:
				var t, err = strconv.ParseInt(*o.Conf["MeterValueSampleInterval"].Value, 10, 64)
				if err != nil {
					o.logger.log(lm, err, assets.Fatal)

					return nil
				}

				sp = sampledEnergyActiveInterval(
					o.Asset.Evses[ie].Connectors[ic].TPowerExport,
					o.Asset.Evses[ie].Connectors[ic].EnergyExport,
					o.Asset.Evses[ie].Connectors[ic].Data[cdp].PowerExport,
					t,
					assets.EnergyActiveExportInterval,
				)

			case assets.EnergyActiveImportInterval:
				var t, err = strconv.ParseInt(*o.Conf["MeterValueSampleInterval"].Value, 10, 64)
				if err != nil {
					o.logger.log(lm, err, assets.Fatal)

					return nil
				}

				sp = sampledEnergyActiveInterval(
					o.Asset.Evses[ie].Connectors[ic].TPower,
					o.Asset.Evses[ie].Connectors[ic].Energy,
					o.Asset.Evses[ie].Connectors[ic].Data[cdp].Power,
					t,
					assets.EnergyActiveImportInterval,
				)

			// This value is static
			case assets.EnergyReactiveExportInterval:
				sp = types.SampledValue{
					Value:     assets.DefReactiveEnergy,
					Unit:      types.UnitOfMeasureVarh,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.EnergyReactiveExportInterval),
					Phase:     types.Phase(assets.Phases[i]),
				}

			// This value is static
			case assets.EnergyReactiveImportInterval:
				sp = types.SampledValue{
					Value:     assets.DefReactiveEnergy,
					Unit:      types.UnitOfMeasureVarh,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.EnergyReactiveImportInterval),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.Frequency:
				sp = types.SampledValue{
					Value:     assets.DefFrequency,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.Frequency),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.PowerActiveExport:
				sp = types.SampledValue{
					Value:     strconv.FormatFloat(float64(o.Asset.Evses[ie].Connectors[ic].Data[cdp].PowerExport), 'f', 4, 64),
					Unit:      types.UnitOfMeasureW,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.PowerActiveExport),
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

			case assets.PowerFactor:
				sp = types.SampledValue{
					Value:     strconv.FormatInt(o.Asset.Evses[ie].Connectors[ic].Data[cdp].PowerFactor, 10),
					Unit:      types.UnitOfMeasurePercent,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.PowerFactor),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.PowerOffered:
				sp = types.SampledValue{
					Value:     strconv.FormatInt(assets.GetConnectorMaxPower(&o.Asset.Evses[ie].Connectors[ic]), 10),
					Unit:      types.UnitOfMeasureW,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.PowerOffered),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.PowerReactiveExport:
				sp = types.SampledValue{
					Value:     assets.DefReactivePower,
					Unit:      types.UnitOfMeasureVar,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.PowerReactiveExport),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.PowerReactiveImport:
				sp = types.SampledValue{
					Value:     assets.DefReactivePower,
					Unit:      types.UnitOfMeasureVar,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.PowerReactiveImport),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.RPM:
				sp = types.SampledValue{
					Value:     assets.DefRPM,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.RPM),
					Phase:     types.Phase(assets.Phases[i]),
				}

			case assets.SoC:
				sp = types.SampledValue{
					Value:     strconv.FormatFloat(o.Asset.Evses[ie].Connectors[ic].CurrentSoC, 'f', 4, 64),
					Unit:      types.UnitOfMeasurePercent,
					Format:    types.ValueFormatRaw,
					Measurand: types.Measurand(assets.SoC),
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
	if o.bootSeq.BootStatus.(core.RegistrationStatus) != core.RegistrationStatusAccepted {
		return
	}

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

	// TODO: this needs to be validated doesnt know if the changes make sense
	var idx, c = getActiveConnector(o.Asset.Evses[0])

	if idx != -1 {
		if c.Data[c.DP.Position].ChargingState == int64(assets.Charging) {
			var cl = strings.Split(*o.Conf["StopTxnAlignedData"].Value, ",")

			var tad []types.MeterValue

			if len(cl) != 0 || cl[0] != "" {
				tad = o.meterValuesAlignedData(cl)
			}

			o.txnAlignedData = append(o.txnAlignedData, tad...)
		}
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
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "meterValuesAlignedData",
		"simulator": o.Asset.Name,
	}

	var spl = []types.SampledValue{}

	var tp = assets.CalculateCPPower(o.Asset.Evses)

	var tpe = assets.CalculateCPPowerExport(o.Asset.Evses)

	for _, conf := range confL {
		var confT = strings.TrimSpace(conf)

		var sp types.SampledValue

		switch confT {
		// TODO: this is only using the evse index 0 to calculate the instantaneous current, since didn't found a way to control multiple evses
		case assets.CurrentImport:
			sp = types.SampledValue{
				Value: strconv.FormatFloat(getAlignedDataCurrent(
					&o.Asset.Evses[0],
					int(o.Asset.Phases),
				), 'f', 4, 64),
				Unit:      types.UnitOfMeasureA,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.CurrentImport),
			}

		// TODO: this is only using the evse index 0 to calculate the max current, since didn't found a way to control multiple evses
		case assets.CurrentOffered:
			sp = types.SampledValue{
				Value:     strconv.FormatFloat(assets.CalculateMaxCurrent(&o.Asset.Evses[0], int(o.Asset.Phases)), 'f', 4, 64),
				Unit:      types.UnitOfMeasureA,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.CurrentOffered),
			}

		case assets.EnergyActiveExportRegister:
			var eet float64

			for _, c := range o.Asset.Evses[0].Connectors {
				eet += c.EnergyExport
			}

			sp = types.SampledValue{
				Value:     strconv.FormatFloat(eet, 'f', 4, 64),
				Unit:      types.UnitOfMeasureWh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyActiveExportRegister),
			}

		case assets.EnergyActiveImportRegister:
			sp = types.SampledValue{
				Value:     strconv.FormatFloat(assets.CalculateCPEnergy(tp, o.st), 'f', 4, 64),
				Unit:      types.UnitOfMeasureWh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyActiveImportRegister),
			}

		// TODO: this is not being calculated and the value is set to 0
		case assets.EnergyReactiveExportRegister:
			sp = types.SampledValue{
				Value:     "0",
				Unit:      types.UnitOfMeasureVarh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyReactiveExportRegister),
			}

		case assets.EnergyReactiveImportRegister:
			sp = types.SampledValue{
				Value:     "0",
				Unit:      types.UnitOfMeasureVarh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyReactiveImportRegister),
			}

		case assets.EnergyActiveExportInterval:
			var t, err = strconv.ParseInt(*o.Conf["ClockAlignedDataInterval"].Value, 10, 64)
			if err != nil {
				o.logger.log(lm, err, assets.Fatal)

				return nil
			}

			sp = alignedEnergyActiveInterval(tpe, t, assets.EnergyActiveExportInterval)

		case assets.EnergyActiveImportInterval:
			var t, err = strconv.ParseInt(*o.Conf["ClockAlignedDataInterval"].Value, 10, 64)
			if err != nil {
				o.logger.log(lm, err, assets.Fatal)

				return nil
			}

			sp = alignedEnergyActiveInterval(tp, t, assets.EnergyActiveImportInterval)

		// This value is static
		case assets.EnergyReactiveExportInterval:
			sp = types.SampledValue{
				Value:     "0",
				Unit:      types.UnitOfMeasureVarh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyReactiveExportInterval),
			}

		// This value is static
		case assets.EnergyReactiveImportInterval:
			sp = types.SampledValue{
				Value:     "0",
				Unit:      types.UnitOfMeasureVarh,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.EnergyReactiveImportInterval),
			}

		case assets.Frequency:
			sp = types.SampledValue{
				Value:     assets.DefFrequency,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.Frequency),
			}

		case assets.PowerActiveExport:
			sp = types.SampledValue{
				Value:     strconv.FormatFloat(tpe, 'f', 4, 64),
				Unit:      types.UnitOfMeasureW,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.PowerActiveExport),
			}

		case assets.PowerActiveImport:
			sp = types.SampledValue{
				Value:     strconv.FormatFloat(tp, 'f', 4, 64),
				Unit:      types.UnitOfMeasureW,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.PowerActiveImport),
			}

		case assets.PowerFactor:
			var tpf int64

			for _, c := range o.Asset.Evses[0].Connectors {
				tpf += c.Data[c.DP.Position].PowerFactor
			}

			sp = types.SampledValue{
				Value:     strconv.FormatInt(tpf, 10),
				Unit:      types.UnitOfMeasurePercent,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.PowerFactor),
			}

		case assets.PowerOffered:
			sp = types.SampledValue{
				Value:     strconv.FormatInt(assets.GetMaxPower(&o.Asset.Evses[0]), 10),
				Unit:      types.UnitOfMeasureW,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.PowerOffered),
			}

		case assets.PowerReactiveExport:
			sp = types.SampledValue{
				Value:     assets.DefReactivePower,
				Unit:      types.UnitOfMeasureVar,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.PowerReactiveExport),
			}

		case assets.PowerReactiveImport:
			sp = types.SampledValue{
				Value:     assets.DefReactivePower,
				Unit:      types.UnitOfMeasureVar,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.PowerReactiveImport),
			}

		case assets.RPM:
			sp = types.SampledValue{
				Value:     assets.DefRPM,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.RPM),
			}

		case assets.SoC:
			var csoc = getAlignedDataSoC(o.Asset.Evses[0])

			sp = types.SampledValue{
				Value:     strconv.FormatInt(int64(csoc), 10),
				Unit:      types.UnitOfMeasurePercent,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.SoC),
			}

		case assets.Voltage:
			var cv = getAlignedDataVoltage(o.Asset.Evses[0])

			sp = types.SampledValue{
				Value:     strconv.FormatInt(cv, 10),
				Unit:      types.UnitOfMeasureV,
				Format:    types.ValueFormatRaw,
				Measurand: types.Measurand(assets.Voltage),
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
// TODO: this function was not updated to remove the loop through the evses list, this func as it is may be relevant to the ocpp 2.0.1
func (o *Ocpp16) updateData() {
	if o.bootSeq.BootStatus.(core.RegistrationStatus) != core.RegistrationStatusAccepted {
		return
	}

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
						o.Asset.Evses[x].Connectors[y].CurrentSoC = calculateSoC(
							o.Asset.Evses[x].Connectors[y].CurrentSoC,
							c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].StartSoC,
							c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].EndSoC,
							c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].Duration,
						)
					} else {
						o.Asset.Evses[x].Connectors[y].DP.Ticker = 0

						if c.DP.Position < int64(len(c.Data)-1) {
							o.Asset.Evses[x].Connectors[y].DP.Position++
						} else {
							o.Asset.Evses[x].Connectors[y].DP.Position = 0
						}

						// reset the current soc when the position changes
						o.Asset.Evses[x].Connectors[y].CurrentSoC =
							o.Asset.Evses[x].Connectors[y].Data[o.Asset.Evses[x].Connectors[y].DP.Position].StartSoC
					}

					// if the charging state has changed in the update
					if cs != c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState {
						o.sendStatusNotification(&o.Asset.Evses[x].Connectors[y])

						// if the connector goes to the finish state
						if c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState == int64(assets.Finishing) {
							o.stopTransaction(o.Asset.Evses[x].CIDTag, &o.Asset.Evses[x].Connectors[y])
							o.txnAlignedData = []types.MeterValue{}
							o.txnSampledData = []types.MeterValue{}

							// if the reset request was used activate the boot sequence and disable the reset trigger
							if o.resetSeq.IsToTrigger {
								o.disconnectSeq = true

								o.connectSeq = true

								o.bootSeq.BootInterval = 0
								o.bootSeq.IsToTrigger = true

								o.resetSeq.IsToTrigger = false
							}
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
											// reset the current soc when the position changes
											o.Asset.Evses[x].Connectors[y].CurrentSoC =
												o.Asset.Evses[x].Connectors[y].Data[o.Asset.Evses[x].Connectors[y].DP.Position].StartSoC

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
					o.Asset.Evses[x].Connectors[y].CurrentSoC = 0
				}
			} else {
				o.notAutoChargePoint(&o.Asset.Evses[x].Connectors[y], x, y)
			}

			o.Asset.Evses[x].Connectors[y].TPower = assets.CalculateTotalPower(c.TPower, c.Data[c.DP.Position].Power)
			o.Asset.Evses[x].Connectors[y].TPowerExport = assets.CalculateTotalPower(c.TPowerExport, c.Data[c.DP.Position].PowerExport)
			o.Asset.Evses[x].Connectors[y].Energy = assets.CalculateEnergy(
				o.Asset.Evses[x].Connectors[y].TPower,
				o.Asset.Evses[x].Connectors[y].Energy,
				c.Data[c.DP.Position].Power,
				o.st,
			)
			o.Asset.Evses[x].Connectors[y].EnergyExport = assets.CalculateEnergy(
				o.Asset.Evses[x].Connectors[y].TPowerExport,
				o.Asset.Evses[x].Connectors[y].EnergyExport,
				c.Data[c.DP.Position].PowerExport,
				o.st,
			)
		}
	}
}

/*
Execute the status notification request with the current charging state of the connector.

c	-	Connector structure with all it's data (*simulator.Connector)
*/
func (o *Ocpp16) sendStatusNotification(c *simulator.Connector) {
	if o.bootSeq.BootStatus.(core.RegistrationStatus) != core.RegistrationStatusAccepted {
		return
	}

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

	var resp, err = o.s.SendRequest(req)

	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
		return
	}

	o.logger.log(lm, resp, assets.Info)
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

	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	// TODO: the transaction message attempts was not tested some investigation needs to be done
	if err != nil {
		o.logger.log(lm, err, assets.Error)

		var tma, errTma = strconv.ParseInt(*o.Conf["TransactionMessageAttempts"].Value, 10, 64)

		if errTma != nil {
			o.logger.log(lm, errTma, assets.Fatal)
		}

		var tmai, errTmai = strconv.ParseInt(*o.Conf["TransactionMessageRetryInterval"].Value, 10, 64)

		if errTmai != nil {
			o.logger.log(lm, errTmai, assets.Fatal)
		}

		for i := 0; i < int(tma); i++ {
			time.Sleep(time.Duration(tmai))

			res, err = o.s.SendRequest(req)

			if err == nil {
				break
			}
		}
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

	var _, err = o.s.SendRequest(req)

	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)

		var tma, errTma = strconv.ParseInt(*o.Conf["TransactionMessageAttempts"].Value, 10, 64)

		if errTma != nil {
			o.logger.log(lm, errTma, assets.Fatal)
		}

		var tmai, errTmai = strconv.ParseInt(*o.Conf["TransactionMessageRetryInterval"].Value, 10, 64)

		if errTmai != nil {
			o.logger.log(lm, errTmai, assets.Fatal)
		}

		for i := 0; i < int(tma); i++ {
			time.Sleep(time.Duration(tmai))

			_, err = o.s.SendRequest(req)

			if err == nil {
				break
			}
		}
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

	o.heartbeatC = 0

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
Logic to handle the heartbeat interval timer, in case reaches the time calls the send heartbeat
function and resets it. If the it didn't reached the time just increments it.
*/
func (o *Ocpp16) processHeartbeat() {
	if o.bootSeq.BootStatus.(core.RegistrationStatus) != core.RegistrationStatusAccepted {
		return
	}

	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "processHeartbeat",
		"feature":   core.HeartbeatFeatureName,
		"simulator": o.Asset.Name,
	}

	o.heartbeatC++

	var hb, hbE = strconv.ParseInt(*o.Conf["HeartbeatInterval"].Value, 10, 64)
	if hbE != nil {
		o.logger.log(lm, hbE, assets.Error)
		return
	}

	if o.heartbeatC == hb {
		o.heartbeat()

		o.heartbeatC = 0
	}
}

// TODO: review the functions that have the send requests name, it might be better to use the same name logic from the ocpp201
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

	// TODO: it is missing the synchronization of the internal clock logic

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

	var res, err = o.s.SendRequest(req)

	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	// TODO: the transaction message attempts was not tested some investigation needs to be done
	if err != nil {
		o.logger.log(lm, err, assets.Error)

		var tma, errTma = strconv.ParseInt(*o.Conf["TransactionMessageAttempts"].Value, 10, 64)

		if errTma != nil {
			o.logger.log(lm, errTma, assets.Fatal)
		}

		var tmai, errTmai = strconv.ParseInt(*o.Conf["TransactionMessageRetryInterval"].Value, 10, 64)

		if errTmai != nil {
			o.logger.log(lm, errTmai, assets.Fatal)
		}

		for i := 0; i < int(tma); i++ {
			time.Sleep(time.Duration(tmai))

			res, err = o.s.SendRequest(req)

			if err == nil {
				break
			}
		}
	}

	o.logger.log(lm, res, assets.Info)
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

	var m, err = strconv.ParseInt(*o.Conf["GetConfigurationMaxKeys"].Value, 10, 64)

	if err != nil {
		var lm2 = map[string]string{
			"protocol":  string(o.Asset.Protocol),
			"function":  "getConfigurationKeys",
			"simulator": o.Asset.Name,
		}

		o.logger.log(lm2, err, assets.Fatal)

		return nil, nil
	}

	for i := range k {
		if i == int(m) {
			break
		}

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
Checks all the connectors from an EVSE and returns false if any of the connectors are enabled.
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
// TODO: this function was not updated to remove the loop through the evses list, this func as it is may be relevant to the ocpp 2.0.1
func (o *Ocpp16) notAutoChargePoint(c *simulator.Connector, x, y int) {
	if c.Enabled {
		if c.DP.Ticker < c.Data[c.DP.Position].Duration {
			o.Asset.Evses[x].Connectors[y].DP.Ticker++
			o.Asset.Evses[x].Connectors[y].CurrentSoC = calculateSoC(
				o.Asset.Evses[x].Connectors[y].CurrentSoC,
				c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].StartSoC,
				c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].EndSoC,
				c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].Duration,
			)
		} else {
			o.Asset.Evses[x].Connectors[y].DP.Ticker = 0

			var cs = c.Data[c.DP.Position].ChargingState

			if assets.Status[c.Data[c.DP.Position].ChargingState] == assets.Status[assets.Finishing] {
				o.Asset.Evses[x].Connectors[y].DP.Position = 0
				o.Asset.Evses[x].Connectors[y].Enabled = false
				o.Asset.Evses[x].Connectors[y].CurrentSoC = 0
				o.txnAlignedData = []types.MeterValue{}
				o.txnSampledData = []types.MeterValue{}

				var aux = o.Asset.Evses[x].Connectors[y]

				switch o.Asset.Evses[x].Connectors[y].Availability {
				case string(assets.Inoperative):
					aux.Data[aux.DP.Position].ChargingState = int64(assets.Unavailable)

				case string(assets.Operative):
					aux.Data[aux.DP.Position].ChargingState = int64(assets.Available)
				}

				go o.sendStatusNotification(&aux)

				return
			}

			// change array data position
			if c.DP.Position < int64(len(c.Data)-1) {
				o.Asset.Evses[x].Connectors[y].DP.Position++
			} else {
				o.Asset.Evses[x].Connectors[y].DP.Position = 0
			}

			// reset the current soc when the position changes
			o.Asset.Evses[x].Connectors[y].CurrentSoC =
				o.Asset.Evses[x].Connectors[y].Data[o.Asset.Evses[x].Connectors[y].DP.Position].StartSoC

			// this loop could cause some issues, this function might need to be reviewed and refactored
			for {
				// break out of the loop in case the reset call was done
				if o.resetSeq.IsToTrigger {
					break
				}

				if c.Data[c.DP.Position].ChargingState == int64(assets.Charging) {
					break
				}

				if c.DP.Position < int64(len(c.Data)-1) {
					o.Asset.Evses[x].Connectors[y].DP.Position++
				} else {
					o.Asset.Evses[x].Connectors[y].DP.Position = 0
				}

				// reset the current soc when the position changes
				o.Asset.Evses[x].Connectors[y].CurrentSoC =
					o.Asset.Evses[x].Connectors[y].Data[o.Asset.Evses[x].Connectors[y].DP.Position].StartSoC
			}

			if cs != c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState {
				o.sendStatusNotification(&o.Asset.Evses[x].Connectors[y])

				if c.Data[o.Asset.Evses[x].Connectors[y].DP.Position].ChargingState == int64(assets.Finishing) {
					o.stopTransaction(o.Asset.Evses[x].CIDTag, &o.Asset.Evses[x].Connectors[y])

					// if the reset request was used, activate the boot sequence and disable the reset trigger
					if o.resetSeq.IsToTrigger {
						o.disconnectSeq = true

						o.connectSeq = true

						o.bootSeq.BootInterval = 0
						o.bootSeq.IsToTrigger = true

						o.resetSeq.IsToTrigger = false
					}
				}
			}
		}
	} else {
		o.Asset.Evses[x].Connectors[y].DP.Position = 0
		o.Asset.Evses[x].Connectors[y].DP.Ticker = 0
		o.Asset.Evses[x].Connectors[y].CurrentSoC = 0
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

/*
Create and return the energy active export and import interval message for the sampled data
(types.SampledValue). This function is created only to not have duplicated code.

tp	-	Total power (float64)

e	-	Energy (float64)

p	-	Current power (int64)

t	-	Time interval (int64)

m	-	Measurand name (string)
*/
func sampledEnergyActiveInterval(tp, e float64, p, t int64, m string) types.SampledValue {
	var sp = types.SampledValue{
		Value: strconv.FormatFloat(assets.CalculateEnergy(
			tp,
			e,
			p,
			time.Now().Add(-time.Second*time.Duration(t)),
		), 'f', 4, 64),
		Unit:      types.UnitOfMeasureWh,
		Format:    types.ValueFormatRaw,
		Measurand: types.Measurand(m),
	}

	return sp
}

/*
Create and return the energy active export and import interval message for the aligned data
(types.SampledValue). This function is created only to not have duplicated code.

p	-	Power of the cp (float64)

t	-	Time interval (int64)

m	-	Measurand name (string)
*/
func alignedEnergyActiveInterval(p float64, t int64, m string) types.SampledValue {
	var sp = types.SampledValue{
		Value:     strconv.FormatFloat(assets.CalculateCPEnergy(p, time.Now().Add(-time.Second*time.Duration(t))), 'f', 4, 64),
		Unit:      types.UnitOfMeasureWh,
		Format:    types.ValueFormatRaw,
		Measurand: types.Measurand(m),
	}

	return sp
}

/*
Calculate the SoC increment to apply, and increment it to the current SoC. It returns the current
SoC result.

c	-	Current SoC

s	-	Data position start SoC

e	-	Data position end SoC

d	-	Data position duration
*/
// TODO: being a percentage the current soc should be converted into int64, not sure if here is the best location to do it
func calculateSoC(c, s, e float64, d int64) float64 {
	var diff = e - s

	var inc = diff / float64(d)

	c += inc

	return c
}

/*
Get the SoC from the active connector. Get the active connector and
return the current SoC for that connector (float64).

e	-	Evse of the asset (simulator.Evse)
*/
func getAlignedDataSoC(e simulator.Evse) float64 {
	var csoc float64

	var i, c = getActiveConnector(e)

	if i != -1 {
		csoc = c.CurrentSoC
	}

	return csoc
}

/*
Get the Voltage from the active connector. Get the active connector, loop the
the phases and return the current voltage for that connector (int64).

e	-	Evse of the asset (simulator.Evse)
*/
func getAlignedDataVoltage(e simulator.Evse) int64 {
	var cv int64

	var i, c = getActiveConnector(e)

	if i != -1 {
		for _, v := range c.Data[c.DP.Position].Voltage {
			if v <= 0 {
				continue
			}

			cv = v
		}
	}

	return cv
}

/*
Get the Evse structure, get the active connector, and calculate the instantaneous current (float64). In
case no active connectors it will return 0.

e	-	The Evse to calculate the current (*simulator.Evse)

ph	-	The asset number of phases (int)
*/
func getAlignedDataCurrent(e *simulator.Evse, ph int) float64 {
	var ci, c = getActiveConnector(*e)

	var cc float64

	if ci > -1 {
		var vi int64

		for _, v := range c.Data[c.DP.Position].Voltage {
			if v > 0 {
				vi = v

				break
			}
		}

		cc = assets.CalculateCurrent(
			c.Data[c.DP.Position].Power,
			c.Data[c.DP.Position].PowerFactor,
			vi,
			int64(ph),
		)
	}

	return cc
}
