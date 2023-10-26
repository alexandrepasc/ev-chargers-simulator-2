package ocpp201

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// TODO: every controller component has an enable variable that could enable or disable the functionality
var components = map[string]component{
	"AlignedDataCtrlr": {
		name:    "AlignedDataCtrlr",
		enabled: true,
	},
	"DeviceDataCtrlr": {
		name:    "DeviceDataCtrlr",
		enabled: true,
		variables: map[string]variable{
			"ItemsPerMessage": {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     "ItemsPerMessage",
							Instance: "GetVariables",
						},
						Component: types.Component{
							Name:     "DeviceDataCtrlr",
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeType:   types.AttributeActual,
						AttributeValue:  "",
					},
					{
						Variable: types.Variable{
							Name:     "ItemsPerMessage",
							Instance: "SetVariables",
						},
						Component: types.Component{
							Name:     "DeviceDataCtrlr",
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeType:   types.AttributeActual,
						AttributeValue:  "",
					},
				},
				mutability: ReadOnly,
			},
		},
	},
	"OCPPCommCtrlr": {
		name:    "OCPPCommCtrlr",
		enabled: true,
		variables: map[string]variable{
			"HeartbeatInterval": {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     "HeartbeatInterval",
							Instance: "",
						},
						Component: types.Component{
							Name:     "OCPPCommCtrlr",
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeValue:  "",
					},
				},
				mutability: ReadWrite,
			},
		},
	},
	"SecurityCtrlr": {
		name:    "SecurityCtrlr",
		enabled: true,
		variables: map[string]variable{
			"BasicAuthPassword": {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     "BasicAuthPassword",
							Instance: "",
						},
						Component: types.Component{
							Name:     "SecurityCtrlr",
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeValue:  "",
					},
				},
				mutability: WriteOnly,
			},
			"Identity": {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     "Identity",
							Instance: "",
						},
						Component: types.Component{
							Name:     "SecurityCtrlr",
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeType:   types.AttributeActual,
						AttributeValue:  "",
					},
				},
				mutability: ReadWrite,
			},
			"SecurityProfile": {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name: "SecurityProfile",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeValue:  "",
					},
				},
				mutability: ReadOnly,
			},
		},
	},
}

type component struct {
	//nolint:structcheck //because dev
	name string
	//nolint:structcheck //because dev
	enabled   bool
	variables map[string]variable
}

type variable struct {
	item       []provisioning.GetVariableResult
	mutability mutability
}

type mutability string

const (
	ReadOnly  mutability = "r"
	WriteOnly mutability = "w"
	ReadWrite mutability = "rw"
)
