package ocpp201

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// TODO: every controller component has an enable variable that could enable or disable the functionality
var components = map[string]component{
	string(alignedDataCtrlr): {
		name:    alignedDataCtrlr,
		enabled: true,
	},
	string(deviceDataCtrlr): {
		name:    deviceDataCtrlr,
		enabled: true,
		variables: map[string]variable{
			string(itemsPerMessage): {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     string(itemsPerMessage),
							Instance: "GetVariables",
						},
						Component: types.Component{
							Name:     string(deviceDataCtrlr),
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeType:   types.AttributeActual,
						AttributeValue:  "",
					},
					{
						Variable: types.Variable{
							Name:     string(itemsPerMessage),
							Instance: "SetVariables",
						},
						Component: types.Component{
							Name:     string(deviceDataCtrlr),
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
	string(oCPPCommCtrlr): {
		name:    oCPPCommCtrlr,
		enabled: true,
		variables: map[string]variable{
			string(heartbeatInterval): {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     string(heartbeatInterval),
							Instance: "",
						},
						Component: types.Component{
							Name:     string(oCPPCommCtrlr),
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeValue:  "",
					},
				},
				mutability: ReadWrite,
			},
			string(networkConfigurationPriority): {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     string(networkConfigurationPriority),
							Instance: "",
						},
						Component: types.Component{
							Name:     string(oCPPCommCtrlr),
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeType:   types.AttributeActual,
						AttributeValue:  "",
					},
				},
				// I'll set this as read only since there is no logic to handle multiple values
				mutability: ReadOnly,
			},
		},
	},
	string(securityCtrlr): {
		name:    securityCtrlr,
		enabled: true,
		variables: map[string]variable{
			string(basicAuthPassword): {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     string(basicAuthPassword),
							Instance: "",
						},
						Component: types.Component{
							Name:     string(securityCtrlr),
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeValue:  "",
					},
				},
				mutability: WriteOnly,
			},
			string(identity): {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name:     string(identity),
							Instance: "",
						},
						Component: types.Component{
							Name:     string(securityCtrlr),
							Instance: "",
						},
						AttributeStatus: provisioning.GetVariableStatusAccepted,
						AttributeType:   types.AttributeActual,
						AttributeValue:  "",
					},
				},
				mutability: ReadWrite,
			},
			string(securityProfile): {
				item: []provisioning.GetVariableResult{
					{
						Variable: types.Variable{
							Name: string(securityProfile),
						},
						Component: types.Component{
							Name:     string(securityCtrlr),
							Instance: "",
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
	name componentName
	//nolint:structcheck //because dev
	enabled   bool
	variables map[string]variable
}

type variable struct {
	item       []provisioning.GetVariableResult
	mutability mutability
}
