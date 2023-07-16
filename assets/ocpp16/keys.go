package ocpp16

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

var config = map[string]core.ConfigurationKey{
	// CORE PROFILE
	// optional boolean
	// "AllowOfflineTxForUnknownId": {
	// 	Key:      "AllowOfflineTxForUnknownId",
	// 	Readonly: false,
	// },
	// optional boolean
	// "AuthorizationCacheEnabled": {
	// 	Key:      "AuthorizationCacheEnabled",
	// 	Readonly: false,
	// },
	// TODO: required and can be readonly boolean, evaluate if we should create a config to set if is read or read and write
	"AuthorizeRemoteTxRequests": {
		Key:      "AuthorizeRemoteTxRequests",
		Readonly: false,
	},
	// optional integer
	// "BlinkRepeat": {
	// 	Key:      "BlinkRepeat",
	// 	Readonly: false,
	// },
	// required integer
	"ClockAlignedDataInterval": {
		Key:      "ClockAlignedDataInterval",
		Readonly: false,
	},
	// required integer this will be hard or impossible to implement since the time out is set on the start of the ocpp server
	"ConnectionTimeOut": {
		Key:      "ConnectionTimeOut",
		Readonly: false,
	},
	// required CSL not sure how to do the logic for this
	// "ConnectorPhaseRotation": {
	// 	Key:      "ConnectorPhaseRotation",
	// 	Readonly: false,
	// },
	// optional integer
	// "ConnectorPhaseRotationMaxLength": {
	// 	Key:      "ConnectorPhaseRotationMaxLength",
	// 	Readonly: true,
	// },
	// required integer
	// "GetConfigurationMaxKeys": {
	// 	Key:      "GetConfigurationMaxKeys",
	// 	Readonly: true,
	// },
	// required integer
	"HeartbeatInterval": {
		Key:      "HeartbeatInterval",
		Readonly: false,
	},
	// optional integer
	// "LightIntensity": {
	// 	Key:      "LightIntensity",
	// 	Readonly: false,
	// },
	// required boolean
	// "LocalAuthorizeOffline": {
	// 	Key:      "LocalAuthorizeOffline",
	// 	Readonly: false,
	// },
	// required boolean
	// "LocalPreAuthorize": {
	// 	Key:      "LocalPreAuthorize",
	// 	Readonly: false,
	// },
	// CSl
	"MeterValuesAlignedData": {
		Key:      "MeterValuesAlignedData",
		Readonly: false,
		Value:    &assets.EnergyActiveImportRegister,
	},
	"MeterValuesSampledData": {
		Key:      "MeterValuesSampledData",
		Readonly: false,
		Value:    &assets.EnergyActiveImportRegister,
	},
	"MeterValueSampleInterval": {
		Key:      "MeterValueSampleInterval",
		Readonly: false,
	},
	// integer
	"NumberOfConnectors": {
		Key:      "NumberOfConnectors",
		Readonly: true,
	},
	// integer
	// "ResetRetries": {
	// 	Key:      "ResetRetries",
	// 	Readonly: false,
	// },
	// boolean
	// "StopTransactionOnEVSideDisconnect": {
	// 	Key:      "StopTransactionOnEVSideDisconnect",
	// 	Readonly: false,
	// },
	// boolean
	// "StopTransactionOnInvalidId": {
	// 	Key:      "StopTransactionOnInvalidId",
	// 	Readonly: false,
	// },
	// CSL
	// "StopTxnAlignedData": {
	// 	Key:      "StopTxnAlignedData",
	// 	Readonly: true,
	// },
	// CSL
	// "StopTxnSampledData": {
	// 	Key:      "StopTxnSampledData",
	// 	Readonly: false,
	// },
	// CSL
	// "SupportedFeatureProfiles": {
	// 	Key:      "SupportedFeatureProfiles",
	// 	Readonly: true,
	// },
	// integer
	// "TransactionMessageAttempts": {
	// 	Key:      "TransactionMessageAttempts",
	// 	Readonly: false,
	// },
	// integer
	// "TransactionMessageRetryInterval": {
	// 	Key:      "TransactionMessageRetryInterval",
	// 	Readonly: false,
	// },
	// boolean
	// "UnlockConnectorOnEVSideDisconnect": {
	// 	Key:      "UnlockConnectorOnEVSideDisconnect",
	// 	Readonly: false,
	// },
	// SMART CHARGING PROFILE
	"ChargingScheduleAllowedChargingRateUnit": {
		Key:      "ChargingScheduleAllowedChargingRateUnit",
		Readonly: false,
		Value:    &assets.Current,
	},
}
