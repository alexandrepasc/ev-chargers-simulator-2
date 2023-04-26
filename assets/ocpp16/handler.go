package ocpp16

import "github.com/lorenzodonini/ocpp-go/ocpp1.6/core"

type Handler struct{}

/**/
func (h *Handler) OnChangeAvailability(_ *core.ChangeAvailabilityRequest) (res *core.ChangeAvailabilityConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnChangeConfiguration(_ *core.ChangeConfigurationRequest) (res *core.ChangeConfigurationConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnClearCache(_ *core.ClearCacheRequest) (res *core.ClearCacheConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnDataTransfer(_ *core.DataTransferRequest) (res *core.DataTransferConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnGetConfiguration(_ *core.GetConfigurationRequest) (res *core.GetConfigurationConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnRemoteStartTransaction(_ *core.RemoteStartTransactionRequest) (res *core.RemoteStartTransactionConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnRemoteStopTransaction(_ *core.RemoteStopTransactionRequest) (res *core.RemoteStopTransactionConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnReset(_ *core.ResetRequest) (res *core.ResetConfirmation, err error) {
	return nil, nil
}

/**/
func (h *Handler) OnUnlockConnector(_ *core.UnlockConnectorRequest) (res *core.UnlockConnectorConfirmation, err error) {
	return nil, nil
}
