package ocpp16

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

type Handler struct{}

/**/
func (o *Ocpp16) OnChangeAvailability(_ *core.ChangeAvailabilityRequest) (res *core.ChangeAvailabilityConfirmation, err error) {
	return res, nil
}

/**/
func (o *Ocpp16) OnChangeConfiguration(_ *core.ChangeConfigurationRequest) (res *core.ChangeConfigurationConfirmation, err error) {
	return res, nil
}

/**/
func (o *Ocpp16) OnClearCache(_ *core.ClearCacheRequest) (res *core.ClearCacheConfirmation, err error) {
	return res, nil
}

/**/
func (o *Ocpp16) OnDataTransfer(_ *core.DataTransferRequest) (res *core.DataTransferConfirmation, err error) {
	return res, nil
}

/**/
func (o *Ocpp16) OnGetConfiguration(req *core.GetConfigurationRequest) (res *core.GetConfigurationConfirmation, err error) {
	o.logger.log(map[string]string{"protocol": string(o.Asset.Protocol), "function": req.GetFeatureName(), "simulator": o.Asset.Name}, req, assets.Info)

	k := o.getConfigurationKeys(req.Key)

	res = &core.GetConfigurationConfirmation{ConfigurationKey: k}

	o.logger.log(map[string]string{"protocol": string(o.Asset.Protocol), "function": req.GetFeatureName(), "simulator": o.Asset.Name}, res, assets.Info)

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
