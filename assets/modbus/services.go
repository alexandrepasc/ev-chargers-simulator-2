package modbus

import (
	"strconv"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/google/uuid"
)

/*
Calculate the power import and return it in W (int).
*/
func (m *Modbus) powerImportW() int {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "powerImportW",
		"simulator": m.Asset.Name,
	}

	var t float64

	var id, err = uuid.Parse("00000000-0000-0000-0000-000000000000")

	if err != nil {
		m.logger.Log(lm, err, assets.Fatal)
	}

	for _, i := range *m.Info {
		if i.UUID == id {
			continue
		}

		t += i.Power
	}

	return int(t)
}

/*
Calculate the power import and return it in kW (int).
*/
func (m *Modbus) powerImportKw() int {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "powerImportKw",
		"simulator": m.Asset.Name,
	}

	var t float64

	var id, err = uuid.Parse("00000000-0000-0000-0000-000000000000")

	if err != nil {
		m.logger.Log(lm, err, assets.Fatal)
	}

	for _, i := range *m.Info {
		if i.UUID == id {
			continue
		}

		t += i.Power
	}

	t /= 1000

	return int(t)
}

/*
Calculate the voltage ampere import and return it in VA (int).
*/
func (m *Modbus) importVa() int {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "importVa",
		"simulator": m.Asset.Name,
	}

	var t float64

	var id, err = uuid.Parse("00000000-0000-0000-0000-000000000000")

	if err != nil {
		m.logger.Log(lm, err, assets.Fatal)
	}

	for _, i := range *m.Info {
		if i.UUID == id {
			continue
		}

		t += i.Power / (float64(i.PowerFactor) / conversion)
	}

	return int(t)
}

/*
Calculate the voltage ampere import and return it in kVA (int).
*/
func (m *Modbus) importKvA() int {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "importKvA",
		"simulator": m.Asset.Name,
	}

	var t float64

	var id, err = uuid.Parse("00000000-0000-0000-0000-000000000000")

	if err != nil {
		m.logger.Log(lm, err, assets.Fatal)
	}

	for _, i := range *m.Info {
		if i.UUID == id {
			continue
		}

		t += (i.Power / conversion) / (float64(i.PowerFactor) / conversion)
	}

	return int(t)
}

/*
Execute the logic to retrieve the value of the specified holding registers address. Will try to
match the string in the models structure with the keys and if it match execute the function. In
case there is no matching key will convert the string to integer. In each case will return the
result (uint16).

reqAddr	-	The request address (uint16)
*/
func (m *Modbus) getHoldingRegistersAddressValue(reqAddr uint16) uint16 {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "getAddressValue",
		"feature":   "HoldingRegisters",
		"simulator": m.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
		"address":   strconv.FormatUint(uint64(reqAddr), 10),
	}

	var ms = m.Mod.Modbus.HoldingRegisters.Addresses[int(reqAddr)]

	var _, ok = functionMap[ms]

	var v uint16

	if ok {
		v = uint16(functionMap[ms].(func(*Modbus) int)(m))
	} else {
		var av, errC = strconv.ParseInt(m.Mod.Modbus.HoldingRegisters.Addresses[int(reqAddr)], 10, 64)

		if errC != nil {
			m.logger.Log(lm, errC, assets.Error)

			return 0
		}

		v = uint16(av)
	}

	return v
}

// func intTo16bitArray(v int) [2]uint16 {
// 	var a16 [2]uint16

// 	a16[0], a16[1] = uint16(v>>bit16), uint16(v)

// 	return a16
// }
