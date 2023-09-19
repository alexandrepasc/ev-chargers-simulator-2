package modbus

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/google/uuid"
)

func (m *Modbus) powerImport() int {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "powerImport",
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

func intTo16bitArray(v int) [2]uint16 {
	var a16 [2]uint16

	a16[0], a16[1] = uint16(v>>bit16), uint16(v)

	return a16
}
