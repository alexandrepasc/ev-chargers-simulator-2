package ocpp201

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// TODO: the correct implementation of the variables could have multiple structures for the same key, review the documentation
var configKeys = map[string]variable{
	"HeartbeatInterval": {
		item: provisioning.GetVariableResult{
			Variable: types.Variable{
				Name: "HeartbeatInterval",
			},
			Component: types.Component{
				Name: "OCPPCommCtrlr",
			},
			AttributeStatus: provisioning.GetVariableStatusAccepted,
			AttributeValue:  "",
		},
		mutability: ReadWrite,
	},
	"ItemsPerMessage": {
		item: provisioning.GetVariableResult{
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
		mutability: ReadOnly,
	},
	"BasicAuthPassword": {
		item: provisioning.GetVariableResult{
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
		mutability: WriteOnly,
	},
	"Identity": {
		item: provisioning.GetVariableResult{
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
		mutability: ReadWrite,
	},
}

type variable struct {
	item       provisioning.GetVariableResult
	mutability mutability
}

type mutability string

const (
	ReadOnly  mutability = "r"
	WriteOnly mutability = "w"
	ReadWrite mutability = "rw"
)
