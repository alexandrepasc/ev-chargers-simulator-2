//nolint:gomnd // because its a draft
package modbus

import (
	"strconv"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/simonvetter/modbus"
)

type Handler struct{}

/*
This method gets called whenever a valid modbus request asking for a coil operation is received
by the server.
*/
func (m *Modbus) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "HandleCoils",
		"feature":   "Coils",
		"simulator": m.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	m.logger.log(lm, req, assets.Info)

	return nil, modbus.ErrIllegalFunction
}

/*
Discrete input handler method.
*/
func (m *Modbus) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "HandleDiscreteInputs",
		"feature":   "DiscreteInputs",
		"simulator": m.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	m.logger.log(lm, req, assets.Info)

	// this is the equivalent of saying
	// "discrete inputs are not supported by this device"
	// (try it with modbus-cli --target tcp://localhost:5502 rdi:1)
	err = modbus.ErrIllegalFunction

	return nil, err
}

/*
This method gets called whenever a valid modbus request asking for a holding register operation
(either read or write) received by the server.
*/
func (m *Modbus) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "HandleHoldingRegisters",
		"feature":   "HoldingRegisters",
		"simulator": m.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	m.logger.log(lm, req, assets.Info)

	switch req.Addr {
	case 2816:
		var ms = m.Mod.Modbus.HoldingRegisters.Addresses[int(req.Addr)]

		var v = intTo16bitArray(functionMap[ms].(func(*Modbus) int)(m))

		for _, vi := range v {
			res = append(res, vi)
		}

	case 3590:
		var ms = m.Mod.Modbus.HoldingRegisters.Addresses[int(req.Addr)]

		var _, ok = functionMap[ms]

		var v [2]uint16

		if ok {
			v = intTo16bitArray(int(functionMap[ms].(func(*Modbus) float64)(m)))
		} else {
			var av, errC = strconv.ParseInt(m.Mod.Modbus.HoldingRegisters.Addresses[int(req.Addr)], 10, 64)

			if errC != nil {
				m.logger.log(lm, errC, assets.Error)

				return nil, modbus.ErrIllegalDataValue
			}

			v = intTo16bitArray(int(av))
		}

		for _, vi := range v {
			res = append(res, vi)
		}

	// any other address is unknown
	default:
		err = modbus.ErrIllegalDataAddress

		m.logger.log(lm, err, assets.Error)

		return nil, err
	}

	lm["sender"] = assets.CP
	lm["type"] = assets.Response

	m.logger.log(lm, res, assets.Info)

	return res, nil
}

/*
This method gets called whenever a valid modbus request asking for an input register operation
is received by the server. Note that input registers are always read-only as per the modbus spec.
*/
func (m *Modbus) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "HandleInputRegisters",
		"feature":   "InputRegisters",
		"simulator": m.Asset.Name,
		"sender":    assets.CS,
		"type":      assets.Request,
	}

	m.logger.log(lm, req, assets.Info)

	return nil, modbus.ErrIllegalFunction
}
