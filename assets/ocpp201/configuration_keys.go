package ocpp201

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// TODO: the correct implementation of the variables could have multiple structures for the same key, review the documentation
var configKeys = map[string]provisioning.GetVariableResult{
	"ItemsPerMessage": {
		Variable: types.Variable{
			Name:     "ItemsPerMessage",
			Instance: "GetVariables",
		},
		Component: types.Component{
			Name: "DeviceDataCtrlr",
		},
		AttributeStatus: provisioning.GetVariableStatusAccepted,
		AttributeType:   types.AttributeActual,
		AttributeValue:  "",
	},
	"BasicAuthPassword": {
		Variable: types.Variable{
			Name: "BasicAuthPassword",
		},
		Component: types.Component{
			Name: "SecurityCtrlr",
		},
		AttributeStatus: provisioning.GetVariableStatusAccepted,
		AttributeType:   types.AttributeActual,
		AttributeValue:  "",
	},
	"Identity": {
		Variable: types.Variable{
			Name: "Identity",
		},
		Component: types.Component{
			Name: "SecurityCtrlr",
		},
		AttributeStatus: provisioning.GetVariableStatusAccepted,
		AttributeType:   types.AttributeActual,
		AttributeValue:  "",
	},
}
