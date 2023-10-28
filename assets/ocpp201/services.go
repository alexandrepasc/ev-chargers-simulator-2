package ocpp201

import (
	"strconv"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

/*
Save the configuration keys to the routine asset, and update the values with the necessery values
defined by the user in the asset and module configuration files.
*/
func (o *Ocpp201) setStartUpConfigurations() {
	o.components = components

	o.connectSeq = true

	o.bootSeq = true

	o.bootReason = provisioning.BootReasonPowerUp

	o.bootStatus = provisioning.RegistrationStatusAccepted

	o.bootInterval = 0

	o.disconnectSeq = false

	o.tick = 0

	o.heartbeatC = 0

	o.secEventSeq = secEventSeq{
		isToTrigger: true,
		eventType:   startupOfTheDevice,
		eventInfo:   "",
	}

	o.logger.log(map[string]string{"protocol": string(o.Asset.Protocol), "function": "setStartUpConfigurations", "simulator": o.Asset.Name},
		o.L.Get(text.StartUpConfigurations), assets.Info)

	var status = provisioning.GetVariableStatusAccepted

	o.components[string(deviceDataCtrlr)].variables[string(itemsPerMessage)] = setComponentVariableValueStatus(
		string(deviceDataCtrlr),
		string(itemsPerMessage),
		strconv.FormatInt(assets.DefGetConfigurationMaxKeys, 10),
		0,
		status,
	)

	o.components[string(deviceDataCtrlr)].variables[string(itemsPerMessage)] = setComponentVariableValueStatus(
		string(deviceDataCtrlr),
		string(itemsPerMessage),
		strconv.FormatInt(assets.DefGetConfigurationMaxKeys, 10),
		1,
		status,
	)

	if !o.Asset.BasicAuth {
		status = provisioning.GetVariableStatusNotSupported
	}

	o.components[string(securityCtrlr)].variables[string(basicAuthPassword)] = setComponentVariableValueStatus(
		string(securityCtrlr),
		string(basicAuthPassword),
		o.Mod.BasicAuth.Password,
		0,
		status,
	)

	o.components[string(securityCtrlr)].variables[string(identity)] = setComponentVariableValueStatus(
		string(securityCtrlr),
		string(identity),
		o.Mod.BasicAuth.Username,
		0,
		status,
	)

	status = provisioning.GetVariableStatusAccepted

	o.components[string(oCPPCommCtrlr)].variables[string(heartbeatInterval)] = setComponentVariableValueStatus(
		string(oCPPCommCtrlr),
		string(heartbeatInterval),
		strconv.FormatInt(assets.DefHeartbeatInterval, 10),
		0,
		status,
	)

	// TODO: this might not be correct
	switch {
	case o.Asset.BasicAuth && o.Asset.TLS:
		o.components[string(securityCtrlr)].variables[string(securityProfile)] = setComponentVariableValueStatus(
			string(securityCtrlr),
			string(securityProfile),
			"3",
			0,
			status,
		)

		o.components[string(oCPPCommCtrlr)].variables[string(networkConfigurationPriority)] = setComponentVariableValueStatus(
			string(oCPPCommCtrlr),
			string(networkConfigurationPriority),
			"3",
			0,
			status,
		)

	case o.Asset.BasicAuth && !o.Asset.TLS:
		o.components[string(securityCtrlr)].variables[string(securityProfile)] = setComponentVariableValueStatus(
			string(securityCtrlr),
			string(securityProfile),
			"1",
			0,
			status,
		)

		o.components[string(oCPPCommCtrlr)].variables[string(networkConfigurationPriority)] = setComponentVariableValueStatus(
			string(oCPPCommCtrlr),
			string(networkConfigurationPriority),
			"1",
			0,
			status,
		)

	case !o.Asset.BasicAuth && !o.Asset.TLS:
		o.components[string(securityCtrlr)].variables[string(securityProfile)] = setComponentVariableValueStatus(
			string(securityCtrlr),
			string(securityProfile),
			"0",
			0,
			status,
		)

		o.components[string(oCPPCommCtrlr)].variables[string(networkConfigurationPriority)] = setComponentVariableValueStatus(
			string(oCPPCommCtrlr),
			string(networkConfigurationPriority),
			"0",
			0,
			status,
		)
	}
}

/*
Sends the boot notification to the central system. Returns the response
(*provisioning.BootNotificationResponse) and the error (error) in case of failure.

r	-	Reason for the boot notification (provisioning.BootReason)
*/
func (o *Ocpp201) sendBootNotification(r provisioning.BootReason) (res *provisioning.BootNotificationResponse, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "sendBootNotification",
		"feature":   provisioning.BootNotificationFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = provisioning.BootNotificationRequest{
		Reason: r,
		ChargingStation: provisioning.ChargingStationType{
			SerialNumber:    o.Mod.Ocpp.SerialNumb,
			Model:           o.Mod.Ocpp.Model,
			VendorName:      o.Mod.Ocpp.Vendor,
			FirmwareVersion: o.Mod.Ocpp.FwVersion,
			Modem: &provisioning.ModemType{
				Iccid: o.Mod.Ocpp.Modem.Iccid,
				Imsi:  o.Mod.Ocpp.Modem.Imsi,
			},
		},
	}

	o.logger.log(lm, req, assets.Info)

	var resp, e = o.s.SendRequest(req)

	// G02.FR.05 reset heartbeat interval when another message has been sent
	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if e != nil {
		o.logger.log(lm, e, assets.Error)
		return nil, e
	}

	o.logger.log(lm, resp.(*provisioning.BootNotificationResponse), assets.Info)

	o.bootStatus = resp.(*provisioning.BootNotificationResponse).Status

	return resp.(*provisioning.BootNotificationResponse), nil
}

/*
Have the logic needed to process the boot notification response.

res	-	The cs boot response (provisioning.BootNotificationResponse)

r	-	Boot reason in case the request needs to be done again (provisioning.BootReason)
*/
func (o *Ocpp201) processBootResponse(res *provisioning.BootNotificationResponse) {
	if res.Status != provisioning.RegistrationStatusAccepted {
		// B03.FR.06 Not accepted and interval grater than 0
		o.bootInterval = res.Interval

		// B03.FR.05 Not accepted and interval is 0
		if res.Interval <= 0 {
			o.bootInterval = int(assets.DefHeartbeatInterval)
		}

		if res.Status == provisioning.RegistrationStatusAccepted {
			o.bootSeq = false
			o.bootInterval = 0
		}
	} else {
		if res.Interval > 0 {
			o.components[string(oCPPCommCtrlr)].variables[string(heartbeatInterval)] = setComponentVariableValueStatus(
				string(oCPPCommCtrlr),
				string(heartbeatInterval),
				strconv.Itoa(res.Interval),
				0,
				provisioning.GetVariableStatusAccepted,
			)
		} else {
			o.components[string(oCPPCommCtrlr)].variables[string(heartbeatInterval)] = setComponentVariableValueStatus(
				string(oCPPCommCtrlr),
				string(heartbeatInterval),
				strconv.FormatInt(assets.DefHeartbeatInterval, 10),
				0,
				provisioning.GetVariableStatusAccepted,
			)
		}

		// B01.FR.06 Synchronization internal clock
		// TODO: need to be done

		for ei, e := range o.Asset.Evses {
			for ci, c := range e.Connectors {
				o.sendStatusNotification(
					e.ID,
					c.ID,
					o.Asset.Evses[ei].Connectors[ci].Data[o.Asset.Evses[ei].Connectors[ci].DP.Position].ChargingState,
				)
			}
		}

		o.bootInterval = 0

		o.bootSeq = false
	}

	o.bootReason = provisioning.BootReasonPowerUp
}

/*
Logic to handle the heartbeat interval timer, in case reaches the time calls the send heartbeat
function and resets it. If the it didn't reached the time just increments it.
*/
func (o *Ocpp201) processHeartbeat() {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "processHeartbeat",
		"feature":   availability.HeartbeatFeatureName,
		"simulator": o.Asset.Name,
	}

	o.heartbeatC++

	var hb, errH = strconv.ParseInt(o.components[string(oCPPCommCtrlr)].variables[string(heartbeatInterval)].item[0].AttributeValue, 10, 64)
	if errH != nil {
		o.logger.log(lm, errH, assets.Error)
		return
	}

	if o.heartbeatC == hb {
		o.sendHeartbeat()

		o.heartbeatC = 0
	}
}

/*
Sends the heartbeat request to the cs.
*/
func (o *Ocpp201) sendHeartbeat() {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "sendHeartbeat",
		"feature":   availability.HeartbeatFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = availability.HeartbeatRequest{}

	o.logger.log(lm, req, assets.Info)

	var resp, e = o.s.SendRequest(req)

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if e != nil {
		o.logger.log(lm, e, assets.Error)
		return
	}

	// G02.FR.06 sync the internal clock
	// G02.FR.07 if heartbeat is never sent send it once every 24 hours

	o.logger.log(lm, resp.(*availability.HeartbeatResponse), assets.Info)
}

/*
Send the connector status information to the cs.

eID	-	EVSE identifier (int64)

cID	-	Connector identifier (int64)

cS	-	Connector status (int64)
*/
func (o *Ocpp201) sendStatusNotification(eID, cID, cS int64) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "sendStatusNotification",
		"feature":   availability.StatusNotificationFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = availability.StatusNotificationRequest{
		Timestamp: &types.DateTime{
			Time: time.Now(),
		},
		EvseID:          int(eID),
		ConnectorID:     int(cID),
		ConnectorStatus: availability.ConnectorStatus(assets.Status[cS]),
	}

	o.logger.log(lm, req, assets.Info)

	var res, err = o.s.SendRequest(req)

	// G02.FR.05 reset heartbeat interval when another message has been sent
	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
		return
	}

	o.logger.log(lm, res.(*availability.StatusNotificationResponse), assets.Info)
}

/*
Get the variables from the asset, validating if they are known or not
([]provisioning.GetVariableResult). If there is any problem with the request returns
the conrresponding error (error).

k	-	The list of variables requested by the cs ([]provisioning.GetVariableData)
*/
func (o *Ocpp201) processGetVariables(k []provisioning.GetVariableData) (r []provisioning.GetVariableResult, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "processGetVariables",
		"feature":   provisioning.GetVariablesFeatureName,
		"simulator": o.Asset.Name,
	}

	var m, errm = strconv.ParseInt(o.components[string(deviceDataCtrlr)].variables[string(itemsPerMessage)].item[0].AttributeValue, 10, 64)
	if errm != nil {
		o.logger.log(lm, errm, assets.Fatal)

		return nil, errm
	}

	// B06.FR.16 More elements than the allowed
	if len(k) > int(m) {
		return nil, ocpp.NewError(ocppj.OccurrenceConstraintViolation, "", "")
	}

	// B06.FR.17 More length than the allowed
	// TODO: not sure how to do this

	for _, ki := range k {
		// B06.FR.06 Unknown component
		var _, okc = components[ki.Component.Name]
		if !okc {
			var uk = provisioning.GetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.GetVariableStatusUnknownComponent,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B06.FR.07 Unknown variable for the given component
		var _, okv = components[ki.Component.Name].variables[ki.Variable.Name]
		if !okv {
			var uk = provisioning.GetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.GetVariableStatusUnknownVariable,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B06.FR.08 Unknown attribute type for the given variable
		var oka = false

		var index int

		for idx, vi := range components[ki.Component.Name].variables[ki.Variable.Name].item {
			if ki.AttributeType == vi.AttributeType {
				oka = true
				index = idx
			}
		}

		if !oka {
			var uk = provisioning.GetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.GetVariableStatusNotSupported,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B06.FR.09 Write only variable
		if components[ki.Component.Name].variables[ki.Variable.Name].mutability == WriteOnly {
			var uk = provisioning.GetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.GetVariableStatusRejected,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B06.FR.11 Empty attribute type
		var okat bool

		switch ki.AttributeType {
		case types.AttributeActual, types.AttributeTarget, types.AttributeMinSet, types.AttributeMaxSet:
			okat = true
		default:
			okat = false
		}

		var vn = provisioning.GetVariableResult{
			Component: components[ki.Component.Name].variables[ki.Variable.Name].item[index].Component,
			Variable:  components[ki.Component.Name].variables[ki.Variable.Name].item[index].Variable,
		}

		if !okat {
			vn.AttributeType = types.AttributeActual
		} else {
			vn.AttributeType = components[ki.Component.Name].variables[ki.Variable.Name].item[index].AttributeType
		}

		// B06.FR.13 No attribute value (at the moment this is created as empty and filled in by the setStartUpConfigurations)

		// B06.FR.14 Instance value provided in the component and/or variable
		// B06.FR.15 No value or empty string instance in the component and/or variable
		if ki.Component.Instance != "" {
			if ki.Component.Instance != components[ki.Component.Name].variables[ki.Variable.Name].item[index].Component.Instance {
				var uk = provisioning.GetVariableResult{
					Variable:        ki.Variable,
					Component:       ki.Component,
					AttributeStatus: provisioning.GetVariableStatusUnknownComponent,
					AttributeType:   ki.AttributeType,
				}

				r = append(r, uk)

				continue
			}
		} else {
			if components[ki.Component.Name].variables[ki.Variable.Name].item[index].Component.Instance != "" {
				var uk = provisioning.GetVariableResult{
					Variable:        ki.Variable,
					Component:       ki.Component,
					AttributeStatus: provisioning.GetVariableStatusUnknownComponent,
					AttributeType:   ki.AttributeType,
				}

				r = append(r, uk)

				continue
			}
		}

		vn.Component.Instance = ki.Component.Instance

		if ki.Variable.Instance != "" {
			if ki.Variable.Instance != components[ki.Component.Name].variables[ki.Variable.Name].item[index].Variable.Instance {
				var uk = provisioning.GetVariableResult{
					Variable:        ki.Variable,
					Component:       ki.Component,
					AttributeStatus: provisioning.GetVariableStatusUnknownVariable,
					AttributeType:   ki.AttributeType,
				}

				r = append(r, uk)

				continue
			}
		} else {
			if components[ki.Component.Name].variables[ki.Variable.Name].item[index].Variable.Instance != "" {
				var uk = provisioning.GetVariableResult{
					Variable:        ki.Variable,
					Component:       ki.Component,
					AttributeStatus: provisioning.GetVariableStatusUnknownVariable,
					AttributeType:   ki.AttributeType,
				}

				r = append(r, uk)

				continue
			}
		}

		vn.Variable.Instance = ki.Variable.Instance

		vn.AttributeStatus = components[ki.Component.Name].variables[ki.Variable.Name].item[index].AttributeStatus
		vn.AttributeValue = components[ki.Component.Name].variables[ki.Variable.Name].item[index].AttributeValue

		r = append(r, vn)
	}

	return r, nil
}

/*
Have all the logic to process the set variable request from the cs. Will do all the checks
and return the response list ([]provisioning.SetVariableResult) and in case some error is found
returns it (error)

k	-	List of requested variables ([]provisioning.SetVariableData)
*/
func (o *Ocpp201) processSetVariables(k []provisioning.SetVariableData) (r []provisioning.SetVariableResult, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "processSetVariables",
		"feature":   provisioning.SetVariablesFeatureName,
		"simulator": o.Asset.Name,
	}

	var m, errm = strconv.ParseInt(o.components[string(deviceDataCtrlr)].variables[string(itemsPerMessage)].item[1].AttributeValue, 10, 64)
	if errm != nil {
		o.logger.log(lm, errm, assets.Fatal)

		return nil, errm
	}

	// B05.FR.11 More elements than the allowed
	if len(k) > int(m) {
		return nil, ocpp.NewError(ocppj.OccurrenceConstraintViolation, "", "")
	}

	// B05.FR.13 Multiple elements with the same component, variable a attribute type combination
	// TODO: not sure how to handle this

	var isBasicAuthPassword = false

	for _, ki := range k {
		// B05.FR.04 Unknown component
		var _, okc = components[ki.Component.Name]
		if !okc {
			var uk = provisioning.SetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.SetVariableStatusUnknownComponent,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B05.FR.05 Unknown variable for given component
		var _, okv = components[ki.Component.Name].variables[ki.Variable.Name]
		if !okv {
			var uk = provisioning.SetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.SetVariableStatusUnknownVariable,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B05.FR.06 Unknow attribute type for given variable
		var oka = false

		var index int

		for idx, vi := range components[ki.Component.Name].variables[ki.Variable.Name].item {
			if ki.AttributeType == vi.AttributeType {
				oka = true
				index = idx
			}
		}

		if !oka {
			var uk = provisioning.SetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.SetVariableStatusNotSupported,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B05.FR.07 Incorrect value format for given variable
		// B05.FR.08 Value is lower or higher than thr variable range

		// B05.FR.09 Read only variable
		if components[ki.Component.Name].variables[ki.Variable.Name].mutability == ReadOnly {
			var uk = provisioning.SetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.SetVariableStatusRejected,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// B05.FR.12 Empty attribute type
		var okat bool

		switch ki.AttributeType {
		case types.AttributeActual, types.AttributeTarget, types.AttributeMinSet, types.AttributeMaxSet:
			okat = true
		default:
			okat = false
		}

		var vn = provisioning.SetVariableResult{
			Component: components[ki.Component.Name].variables[ki.Variable.Name].item[index].Component,
			Variable:  components[ki.Component.Name].variables[ki.Variable.Name].item[index].Variable,
		}

		if !okat {
			vn.AttributeType = types.AttributeActual
		} else {
			vn.AttributeType = components[ki.Component.Name].variables[ki.Variable.Name].item[index].AttributeType
		}

		vn.Variable.Instance = components[ki.Component.Name].variables[ki.Variable.Name].item[index].Variable.Instance

		vn.AttributeStatus = provisioning.SetVariableStatusAccepted
		components[ki.Component.Name].variables[ki.Variable.Name].item[index].AttributeValue = ki.AttributeValue

		if ki.Variable.Name == string(basicAuthPassword) {
			isBasicAuthPassword = true
		}

		r = append(r, vn)
	}

	// A01 Update cp password for http basic authentication
	if isBasicAuthPassword {
		o.disconnectSeq = true

		o.connectSeq = true

		o.bootSeq = true

		o.bootReason = provisioning.BootReasonApplicationReset

		o.secEventSeq = secEventSeq{
			isToTrigger: true,
			eventType:   resetOrReboot,
			eventInfo:   "",
		}
	}

	return r, nil
}

/*
Send the security event notification to the cs.

t	-	The security event type from the list of events (securityEventType)

i	-	Additional information to add to the request (string)
*/
func (o *Ocpp201) sendSecurityEventNotification(t securityEventType, i string) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "sendSecurityEventNotification",
		"feature":   security.SecurityEventNotificationFeatureName,
		"simulator": o.Asset.Name,
		"sender":    assets.CP,
		"type":      assets.Request,
	}

	var req = security.SecurityEventNotificationRequest{
		Type: string(t),
		Timestamp: &types.DateTime{
			Time: time.Now(),
		},
		TechInfo: i,
	}

	o.logger.log(lm, req, assets.Info)

	var res, err = o.s.SendRequest(req)

	// G02.FR.05 reset heartbeat interval when another message has been sent
	o.heartbeatC = 0

	lm["sender"] = assets.CS
	lm["type"] = assets.Response

	if err != nil {
		o.logger.log(lm, err, assets.Error)
		return
	}

	o.logger.log(lm, res.(*security.SecurityEventNotificationResponse), assets.Info)

	o.secEventSeq.isToTrigger = false
}

/*
B03.FR.08 Boot rejected and not trigger message BootNotification

Checks if the boot status is rejected, if so returns the ocpp security error (error),
if not returns nil.

lm	-	The logging fields (map[string]string)
*/
func (o *Ocpp201) isBootRejected(lm map[string]string) error {
	if o.bootStatus == provisioning.RegistrationStatusRejected {
		o.logger.log(lm, ocppj.SecurityError, assets.Error)

		return ocpp.NewError(ocppj.SecurityError, "", "")
	}

	return nil
}

/*
Set the variable value and status sent in the parameters, and return the variable structure
with the list of variables (variable).

cn	-	Component name where the variable lives (string)

vn	-	Name of the variable to change (string)

val	-	Value to set to the variable (string)

i	-	Index of the variable instance (int)

s	-	Status to set to the variable (provisioning.GetVariableStatus)
*/
func setComponentVariableValueStatus(cn, vn, val string, i int, s provisioning.GetVariableStatus) variable {
	var uis = []provisioning.GetVariableResult{}

	for x, y := range components[cn].variables[vn].item {
		if x == i {
			var ui = provisioning.GetVariableResult{
				Variable:        components[cn].variables[vn].item[x].Variable,
				Component:       components[cn].variables[vn].item[x].Component,
				AttributeStatus: s,
				AttributeType:   components[cn].variables[vn].item[x].AttributeType,
				AttributeValue:  val,
			}

			uis = append(uis, ui)

			continue
		}

		uis = append(uis, y)
	}

	return variable{
		item:       uis,
		mutability: components[cn].variables[vn].mutability,
	}
}
