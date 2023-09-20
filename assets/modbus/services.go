package modbus

import (
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
		m.logger.log(lm, err, assets.Fatal)
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
		m.logger.log(lm, err, assets.Fatal)
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

// func intTo16bitArray(v int) [2]uint16 {
// 	var a16 [2]uint16

// 	a16[0], a16[1] = uint16(v>>bit16), uint16(v)

// 	return a16
// }
