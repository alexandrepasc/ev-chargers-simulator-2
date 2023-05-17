package ocpp16

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
)

type Handler struct{}

/**/
func (o *Ocpp16) OnChangeAvailability(_ *core.ChangeAvailabilityRequest) (res *core.ChangeAvailabilityConfirmation, err error) {
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
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	k, u := o.getConfigurationKeys(req.Key)

	res = &core.GetConfigurationConfirmation{ConfigurationKey: k}

	res.UnknownKey = u

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
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &core.ChangeConfigurationConfirmation{Status: o.setConfiguration(req)}

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
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	o.Auth = nil

	res = &core.ClearCacheConfirmation{Status: core.ClearCacheStatusAccepted}

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
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	res = &core.DataTransferConfirmation{Status: core.DataTransferStatusAccepted}

	if req.VendorId != o.Mod.Ocpp.VendorID {
		res = &core.DataTransferConfirmation{Status: core.DataTransferStatusUnknownVendorId}
	}

	lm["type"] = assets.Response

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

/**/
func (o *Ocpp16) OnRemoteStartTransaction(_ *core.RemoteStartTransactionRequest) (res *core.RemoteStartTransactionConfirmation, err error) {
	return res, nil
}

/**/
func (o *Ocpp16) OnRemoteStopTransaction(_ *core.RemoteStopTransactionRequest) (res *core.RemoteStopTransactionConfirmation, err error) {
	return res, nil
}

/**/
func (o *Ocpp16) OnReset(_ *core.ResetRequest) (res *core.ResetConfirmation, err error) {
	return res, nil
}

/**/
func (o *Ocpp16) OnUnlockConnector(_ *core.UnlockConnectorRequest) (res *core.UnlockConnectorConfirmation, err error) {
	return res, nil
}

func (o *Ocpp16) OnGetLocalListVersion(_ *localauth.GetLocalListVersionRequest) (res *localauth.GetLocalListVersionConfirmation, err error) {
	return res, nil
}

func (o *Ocpp16) OnSendLocalList(_ *localauth.SendLocalListRequest) (res *localauth.SendLocalListConfirmation, err error) {
	return res, nil
}
