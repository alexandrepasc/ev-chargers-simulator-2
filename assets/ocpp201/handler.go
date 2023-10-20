//nolint:gocritic,revive,whitespace,wsl //because dev
package ocpp201

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

func (o *Ocpp201) OnGetVariables(req *provisioning.GetVariablesRequest) (res *provisioning.GetVariablesResponse, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnGetVariables",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	var r, e = o.processGetVariables(req.GetVariableData)
	if e != nil {
		o.logger.log(lm, e, assets.Error)

		return nil, e
	}

	res = &provisioning.GetVariablesResponse{
		GetVariableResult: r,
	}

	o.logger.log(lm, res, assets.Info)

	return res, nil
}

func (o *Ocpp201) OnSetVariables(req *provisioning.SetVariablesRequest) (res *provisioning.SetVariablesResponse, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnSetVariables",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.log(lm, req, assets.Info)
	// logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnGetBaseReport(req *provisioning.GetBaseReportRequest) (res *provisioning.GetBaseReportResponse, err error) {
	// logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnGetReport(req *provisioning.GetReportRequest) (res *provisioning.GetReportResponse, err error) {
	// logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnReset(req *provisioning.ResetRequest) (res *provisioning.ResetResponse, err error) {
	// logDefault(request.GetFeatureName()).Info("reset handled")
	res = provisioning.NewResetResponse(provisioning.ResetStatusAccepted)
	return
}

func (o *Ocpp201) OnSetNetworkProfile(req *provisioning.SetNetworkProfileRequest) (res *provisioning.SetNetworkProfileResponse, err error) {
	// logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnRequestStartTransaction(req *remotecontrol.RequestStartTransactionRequest) (res *remotecontrol.RequestStartTransactionResponse, err error) {

	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnRequestStopTransaction(req *remotecontrol.RequestStopTransactionRequest) (res *remotecontrol.RequestStopTransactionResponse, err error) {

	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnTriggerMessage(req *remotecontrol.TriggerMessageRequest) (res *remotecontrol.TriggerMessageResponse, err error) {
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnUnlockConnector(req *remotecontrol.UnlockConnectorRequest) (res *remotecontrol.UnlockConnectorResponse, err error) {
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}
