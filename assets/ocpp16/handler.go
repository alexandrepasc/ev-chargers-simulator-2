package ocpp16

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type Handler struct{}

/*
Process the change availability request changing the connector status from operative to inoperative
or the other way around. In case a session is occurring the connector it will change the state after
the session finishes.
*/
func (o *Ocpp16) OnChangeAvailability(req *core.ChangeAvailabilityRequest) (res *core.ChangeAvailabilityConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnChangeAvailability",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = o.processChangeAvailability(req)

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Get the central system get configuration request and sends the response to the central system the
response with the requested keys (*core.GetConfigurationConfirmation, error).

req	-	Get configuration request (*core.GetConfigurationRequest)
*/
func (o *Ocpp16) OnGetConfiguration(req *core.GetConfigurationRequest) (res *core.GetConfigurationConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnGetConfiguration",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	k, u := o.getConfigurationKeys(req.Key)

	res = &core.GetConfigurationConfirmation{ConfigurationKey: k}

	res.UnknownKey = u

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Get the change configuration request from the CS and prosses it. The request will send a key value
pair and in case the key matches the listed configurations the value will be changed.
*/
func (o *Ocpp16) OnChangeConfiguration(req *core.ChangeConfigurationRequest) (res *core.ChangeConfigurationConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnChangeConfiguration",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &core.ChangeConfigurationConfirmation{Status: o.setConfiguration(req)}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Receives the clear cache request and proceed to remove all the data stored related from the
authorization list.
*/
func (o *Ocpp16) OnClearCache(req *core.ClearCacheRequest) (res *core.ClearCacheConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnClearCache",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	o.localAuth.version = 0
	o.localAuth.list = nil

	res = &core.ClearCacheConfirmation{Status: core.ClearCacheStatusAccepted}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Receives validate the vendor id and logs the request and the response.
*/
func (o *Ocpp16) OnDataTransfer(req *core.DataTransferRequest) (res *core.DataTransferConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnDataTransfer",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &core.DataTransferConfirmation{Status: core.DataTransferStatusAccepted}

	if req.VendorId != o.Mod.Ocpp.VendorID {
		res = &core.DataTransferConfirmation{Status: core.DataTransferStatusUnknownVendorId}
	}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Handles the start transaction request from the CS, process the information and start a charging
session.
*/
func (o *Ocpp16) OnRemoteStartTransaction(req *core.RemoteStartTransactionRequest) (res *core.RemoteStartTransactionConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnRemoteStartTransaction",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}

	if canEnable(o.Asset.Evses[0].Connectors) {
		res = &core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}

		go o.processRemoteStartTransaction(req)
	}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Handles the CS stop transaction request and process it and stop the transaction in case one is
occurring in the connector went in the request.
*/
func (o *Ocpp16) OnRemoteStopTransaction(req *core.RemoteStopTransactionRequest) (res *core.RemoteStopTransactionConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnRemoteStopTransaction",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = o.processRemoteStopTransaction(req)

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Receives the reset request from the CS and make the required logic in each of the reset cases
(hard, soft).
*/
func (o *Ocpp16) OnReset(req *core.ResetRequest) (res *core.ResetConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnReset",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = o.processReset(req)

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Handles the unlock connector request from the CS, this is used by te support team in case the
client can't remove the connector from the EV. This validates that the connector ID matches with
any of the connectors and responds unlocked, in case it doesn't match returns unlock failed.
*/
func (o *Ocpp16) OnUnlockConnector(req *core.UnlockConnectorRequest) (res *core.UnlockConnectorConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnUnlockConnector",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &core.UnlockConnectorConfirmation{Status: core.UnlockStatusUnlockFailed}

	// at the moment this is only supporting 1 evse for simulator so the evse is set to 0
	for _, c := range o.Asset.Evses[0].Connectors {
		if c.ID != int64(req.ConnectorId) {
			continue
		}

		res = &core.UnlockConnectorConfirmation{Status: core.UnlockStatusUnlocked}
	}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Receives the CS request and returns the local list version.
*/
func (o *Ocpp16) OnGetLocalListVersion(req *localauth.GetLocalListVersionRequest) (res *localauth.GetLocalListVersionConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnGetLocalListVersion",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &localauth.GetLocalListVersionConfirmation{ListVersion: int(o.localAuth.version)}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
Handles the request filter if the feature is supported. In case it is supported execute the logic
to update the local list.
*/
func (o *Ocpp16) OnSendLocalList(req *localauth.SendLocalListRequest) (res *localauth.SendLocalListConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnSendLocalList",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	if o.Asset.AuthList {
		res = o.processSendLocalList(req)
	} else {
		res = &localauth.SendLocalListConfirmation{Status: localauth.UpdateStatusNotSupported}
	}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

func (o *Ocpp16) OnTriggerMessage(req *remotetrigger.TriggerMessageRequest) (res *remotetrigger.TriggerMessageConfirmation, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnTriggerMessage",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &remotetrigger.TriggerMessageConfirmation{Status: remotetrigger.TriggerMessageStatusNotImplemented}

	return res, nil
}
