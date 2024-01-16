//nolint:revive //because dev
package ocpp201

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

/*
Handles the get variables request from the cs and returns the response or error in case of failure.
*/
//nolint:dupl //because at the moment there is no idea to clean this
func (o *Ocpp201) OnGetVariables(req *provisioning.GetVariablesRequest) (res *provisioning.GetVariablesResponse, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnGetVariables",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.Log(lm, req, assets.Info)

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	// B03.FR.08 Boot rejected and not trigger message BootNotification
	var eB = o.isBootRejected(lm)
	if eB != nil {
		return nil, ocpp.NewError(ocppj.SecurityError, "", "")
	}

	var r, e = o.processGetVariables(req.GetVariableData)
	if e != nil {
		o.logger.Log(lm, e, assets.Error)

		return nil, e
	}

	res = &provisioning.GetVariablesResponse{
		GetVariableResult: r,
	}

	o.logger.Log(lm, res, assets.Info)

	return res, nil
}

/*
Handles the set variables request from the cs and returns the response or error in case of failure.
*/
//nolint:dupl //because at the moment there is no idea to clean this
func (o *Ocpp201) OnSetVariables(req *provisioning.SetVariablesRequest) (res *provisioning.SetVariablesResponse, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnSetVariables",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.Log(lm, req, assets.Info)

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	// B03.FR.08 Boot rejected and not trigger message BootNotification
	var eB = o.isBootRejected(lm)
	if eB != nil {
		return nil, ocpp.NewError(ocppj.SecurityError, "", "")
	}

	var r, e = o.processSetVariables(req.SetVariableData)
	if e != nil {
		o.logger.Log(lm, e, assets.Error)

		return nil, e
	}

	res = &provisioning.SetVariablesResponse{
		SetVariableResult: r,
	}

	o.logger.Log(lm, res, assets.Info)

	return res, nil
}

/*
Handles the trigger message request from the cs and returns the response or error in case of
failure.
*/
func (o *Ocpp201) OnTriggerMessage(req *remotecontrol.TriggerMessageRequest) (res *remotecontrol.TriggerMessageResponse, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnTriggerMessage",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.Log(lm, req, assets.Info)
	// B02.FR.02
	// B02.FR.09 Boot pending returns the boot response
	// B03.FR.08 Boot rejected and not trigger message BootNotification

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	//nolint:exhaustive //because still in dev
	switch req.RequestedMessage {
	case remotecontrol.MessageTriggerBootNotification:
		o.bootSeq.bootReason = provisioning.BootReasonTriggered
		o.bootSeq.isToTrigger = true

		return &remotecontrol.TriggerMessageResponse{Status: remotecontrol.TriggerMessageStatusAccepted}, nil
	default:
		return nil, ocpp.NewError(ocppj.NotSupported, "", "")
	}
}

func (o *Ocpp201) OnReset(req *provisioning.ResetRequest) (res *provisioning.ResetResponse, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "OnReset",
		"feature":   req.GetFeatureName(),
		"simulator": o.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	o.logger.Log(lm, req, assets.Info)

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	// B03.FR.08 Boot rejected returns SecurityError
	var eB = o.isBootRejected(lm)
	if eB != nil {
		return nil, ocpp.NewError(ocppj.SecurityError, "", "")
	}

	var ty, inf = o.processResetRequest(req.EvseID, req.Type)

	res = &provisioning.ResetResponse{
		Status:     ty,
		StatusInfo: inf,
	}

	o.logger.Log(lm, res, assets.Info)

	return res, nil
}

func (o *Ocpp201) OnGetBaseReport(req *provisioning.GetBaseReportRequest) (res *provisioning.GetBaseReportResponse, err error) {
	// B02.FR.02
	// B02.FR.09 Boot pending returns the get base report response
	// B03.FR.08 Boot rejected returns SecurityError
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnGetReport(req *provisioning.GetReportRequest) (res *provisioning.GetReportResponse, err error) {
	// B02.FR.02
	// B02.FR.09 Boot pending returns the get report response
	// B03.FR.08 Boot rejected returns SecurityError
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnSetNetworkProfile(req *provisioning.SetNetworkProfileRequest) (res *provisioning.SetNetworkProfileResponse, err error) {
	// B03.FR.08 Boot rejected returns SecurityError
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnRequestStartTransaction(req *remotecontrol.RequestStartTransactionRequest) (res *remotecontrol.RequestStartTransactionResponse, err error) {
	// B02.FR.05 Boot pending returns the rejected response
	// B03.FR.08 Boot rejected returns SecurityError
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnRequestStopTransaction(req *remotecontrol.RequestStopTransactionRequest) (res *remotecontrol.RequestStopTransactionResponse, err error) {
	// B02.FR.05 Boot pending returns the rejected response
	// B03.FR.08 Boot rejected returns SecurityError
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}

func (o *Ocpp201) OnUnlockConnector(req *remotecontrol.UnlockConnectorRequest) (res *remotecontrol.UnlockConnectorResponse, err error) {
	// B03.FR.08 Boot rejected returns SecurityError
	return nil, ocpp.NewError(ocppj.NotSupported, "Not supported", "")
}
