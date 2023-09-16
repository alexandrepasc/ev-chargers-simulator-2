//nolint:gomnd // because its a draft
package modbus

import (
	"fmt"

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
	var regAddr uint16

	if req.UnitId != 1 {
		// only accept unit ID #1
		err = modbus.ErrIllegalFunction
		return nil, err
	}

	fmt.Println("___________________________________________________________")
	fmt.Println("HandleHoldingRegisters")
	fmt.Println(req)

	// since we're manipulating variables shared between multiple goroutines,
	// acquire a lock to avoid concurrency issues.
	m.lock.Lock()
	// release the lock upon return
	defer m.lock.Unlock()

	// loop through `quantity` registers
	for i := 0; i < int(req.Quantity); i++ {
		// compute the target register address
		regAddr = req.Addr + uint16(i)

		switch regAddr {
		// expose the static, read-only value of 0xff00 in register 100
		case 100:
			res = append(res, 0xff00)

		// expose holdingReg1 in register 101 (RW)
		case 101:
			if req.IsWrite {
				m.holdingReg1 = req.Args[i]
			}

			res = append(res, m.holdingReg1)

		// expose holdingReg2 in register 102 (RW)
		case 102:
			if req.IsWrite {
				// only accept values 2 and 4
				switch req.Args[i] {
				case 2, 4:
					m.holdingReg2 = req.Args[i]

					// make note of the change (e.g. for auditing purposes)
					fmt.Printf("%s set reg#102 to %v\n", req.ClientAddr, m.holdingReg2)
				default:
					// if the written value is neither 2 nor 4,
					// return a modbus "illegal data value" to
					// let the client know that the value is
					// not acceptable.
					err = modbus.ErrIllegalDataValue
					return nil, err
				}
			}

			res = append(res, m.holdingReg2)

		// expose eh.holdingReg3 in register 103 (RW)
		// note: eh.holdingReg3 is a signed 16-bit integer
		case 103:
			if req.IsWrite {
				// cast the 16-bit unsigned integer passed by the server
				// to a 16-bit signed integer when writing
				m.holdingReg3 = int16(req.Args[i])
			}
			// cast the 16-bit signed integer from the handler to a 16-bit unsigned
			// integer so that we can append it to `res`.
			res = append(res, uint16(m.holdingReg3))

		// expose the 16 most-significant bits of eh.holdingReg4 in register 200
		case 200:
			if req.IsWrite {
				m.holdingReg4 =
					((uint32(req.Args[i])<<16)&0xffff0000 |
						(m.holdingReg4 & 0x0000ffff))
			}

			res = append(res, uint16((m.holdingReg4>>16)&0x0000ffff))

		// expose the 16 least-significant bits of eh.holdingReg4 in register 201
		case 201:
			if req.IsWrite {
				m.holdingReg4 =
					(uint32(req.Args[i])&0x0000ffff |
						(m.holdingReg4 & 0xffff0000))
			}

			res = append(res, uint16(m.holdingReg4&0x0000ffff))

		// any other address is unknown
		default:
			err = modbus.ErrIllegalDataAddress
			return nil, err
		}
	}
	fmt.Println("___________________________________________________________")
	fmt.Println(res)

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
