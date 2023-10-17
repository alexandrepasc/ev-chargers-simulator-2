package ocpp201

import (
	"strconv"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

/*
Save the configuration keys to the routine asset, and update the values with the necessery values
defined by the user in the asset and module configuration files.
*/
func (o *Ocpp201) setStartUpConfigurations() {
	o.ConfigKeys = configKeys

	o.logger.log(map[string]string{"protocol": string(o.Asset.Protocol), "function": "setStartUpConfigurations", "simulator": o.Asset.Name},
		o.L.Get(text.StartUpConfigurations), assets.Info)

	var status = provisioning.GetVariableStatusAccepted

	o.ConfigKeys["ItemsPerMessage"] = getConfigKey("ItemsPerMessage", strconv.FormatInt(assets.DefGetConfigurationMaxKeys, 10), status)

	if !o.Asset.BasicAuth {
		status = provisioning.GetVariableStatusNotSupported
	}

	o.ConfigKeys["BasicAuthPassword"] = getConfigKey("BasicAuthPassword", o.Mod.BasicAuth.Password, status)

	o.ConfigKeys["Identity"] = getConfigKey("Identity", o.Mod.BasicAuth.Username, status)
}

/*
Get the configuration keys from the asset, validating if they are known or not
([]provisioning.GetVariableResult). If there is any problem with the request returns
the conrresponding error (error).

k	-	The list of variables requested by the cs ([]provisioning.GetVariableData)
*/
// TODO: validate that the attribute type matches
// TODO: validate the component value
func (o *Ocpp201) getConfigurationKeys(k []provisioning.GetVariableData) (r []provisioning.GetVariableResult, err error) {
	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "getConfigurationKeys",
		"feature":   "GetVariables",
		"simulator": o.Asset.Name,
	}

	var m, errm = strconv.ParseInt(o.ConfigKeys["ItemsPerMessage"].AttributeValue, 10, 64)
	if errm != nil {
		o.logger.log(lm, errm, assets.Fatal)

		return nil, errm
	}

	if len(k) > int(m) {
		return nil, ocpp.NewError(ocppj.OccurrenceConstraintViolation, "", "")
	}

	for _, ki := range k {
		var _, ok = o.ConfigKeys[ki.Variable.Name]

		// The variable is not listed in the configuration keys
		if !ok {
			var uk = provisioning.GetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.GetVariableStatusUnknownVariable,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// The requested variable component is not equal to the configuration key
		if o.ConfigKeys[ki.Variable.Name].Component.Name != ki.Component.Name || o.ConfigKeys[ki.Variable.Name].Component.Instance != ki.Component.Instance {
			var uk = provisioning.GetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.GetVariableStatusUnknownComponent,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		// The requested attribute type is not equal to the configuraion key
		if o.ConfigKeys[ki.Variable.Name].AttributeType != ki.AttributeType {
			var uk = provisioning.GetVariableResult{
				Variable:        ki.Variable,
				Component:       ki.Component,
				AttributeStatus: provisioning.GetVariableStatusNotSupported,
				AttributeType:   ki.AttributeType,
			}

			r = append(r, uk)

			continue
		}

		r = append(r, o.ConfigKeys[ki.Variable.Name])
	}

	return r, nil
}

/*
Ge the configuration name, the value, the status and return the structure to be set with the
correct values (provisioning.GetVariableResult).

n	-	The configuration variable name (string)

v	-	The value to set in the configuration (string)

s	-	The status of the configuration variable (provisioning.GetVariableStatus)
*/
func getConfigKey(n, v string, s provisioning.GetVariableStatus) provisioning.GetVariableResult {
	return provisioning.GetVariableResult{
		Variable:        configKeys[n].Variable,
		Component:       configKeys[n].Component,
		AttributeStatus: s,
		AttributeType:   configKeys[n].AttributeType,
		AttributeValue:  v,
	}
}
