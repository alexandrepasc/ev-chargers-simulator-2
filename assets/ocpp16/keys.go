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
	// no logic added to this configuration
	"ConnectorPhaseRotation": {
		Key:      "ConnectorPhaseRotation",
		Readonly: false,
	},
	// optional integer
	// "ConnectorPhaseRotationMaxLength": {
	// 	Key:      "ConnectorPhaseRotationMaxLength",
	// 	Readonly: true,
	// },
	"GetConfigurationMaxKeys": {
		Key:      "GetConfigurationMaxKeys",
		Readonly: true,
	},
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
	// this key doesn't has any logic associated, since the ideia of this project is to be used with a cs
	"LocalAuthorizeOffline": {
		Key:      "LocalAuthorizeOffline",
		Readonly: false,
	},
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
	"NumberOfConnectors": {
		Key:      "NumberOfConnectors",
		Readonly: true,
	},
	// TODO: this value is static hand does not have any logic
	"ResetRetries": {
		Key:      "ResetRetries",
		Readonly: false,
	},
	// TODO: this value is static hand does not have any logic
	"StopTransactionOnEVSideDisconnect": {
		Key:      "StopTransactionOnEVSideDisconnect",
		Readonly: false,
	},
	"StopTransactionOnInvalidId": {
		Key:      "StopTransactionOnInvalidId",
		Readonly: false,
	},
	"StopTxnAlignedData": {
		Key:      "StopTxnAlignedData",
		Readonly: false,
	},
	"StopTxnSampledData": {
		Key:      "StopTxnSampledData",
		Readonly: false,
	},
	"SupportedFeatureProfiles": {
		Key:      "SupportedFeatureProfiles",
		Readonly: true,
		Value:    assets.GetStringPointer(assets.DefSupportedFeatureProfiles),
	},
	"TransactionMessageAttempts": {
		Key:      "TransactionMessageAttempts",
		Readonly: false,
	},
	"TransactionMessageRetryInterval": {
		Key:      "TransactionMessageRetryInterval",
		Readonly: false,
	},
	// no logic added to this configuration
	"UnlockConnectorOnEVSideDisconnect": {
		Key:      "UnlockConnectorOnEVSideDisconnect",
		Readonly: false,
	},
	// SMART CHARGING PROFILE
	"ChargingScheduleAllowedChargingRateUnit": {
		Key:      "ChargingScheduleAllowedChargingRateUnit",
		Readonly: false,
		Value:    &assets.Current,
	},
}
