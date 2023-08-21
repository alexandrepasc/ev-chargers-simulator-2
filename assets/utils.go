package assets

import (
	"math"
	"strconv"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
)

/*
Calculate the current from power and voltage, and in case of 3 phases the power factor will be
used.

It will return the result in amps (float64).

p	-	Power in W (int64)

pf	-	Power factor (int64)

v	-	Voltage in V (int64)

ph	-	Phases number (int64)
*/
// TODO: need to have in attemption the current type to calculate this
func CalculateCurrent(p, pf, v, ph int64) float64 {
	const (
		sqrt  float64 = 3
		pfper float64 = 1000
	)

	if ph == 1 {
		var t = float64(p) / float64(v)

		return t
	}

	var t = float64(p) / (math.Sqrt(sqrt) * (float64(pf) / pfper) * float64(v))

	return t
}

/*
pp	-	Previous power stored in the connector (float64)

cp	-	Current connector power (float64)

st	-	Simulator start time (time.Time)
*/
func CalculateTotalPower(pp float64, cp int64) (tp float64) {
	const med float64 = 2

	if cp == 0 {
		return pp
	}

	tp = (pp + float64(cp)) / med

	return tp
}

/**/
func CalculateEnergy(tp, ce float64, cp int64, st time.Time) (e float64) {
	if cp == 0 {
		return ce
	}

	e = tp * time.Since(st).Hours()

	return e
}

/*
Sum the power of all the connectors of the CP and returns the calculated value (float64).

el	-	The list of evses that the CP has ([]simulator.Evse)
*/
func CalculateCPPower(el []simulator.Evse) float64 {
	var t float64

	for _, e := range el {
		for _, c := range e.Connectors {
			t += c.TPower
		}
	}

	return t
}

/*
Sum the power export of all the connectors of the CP and returns the calculated value (float64).

el	-	The list of evses that the CP has ([]simulator.Evse)
*/
func CalculateCPPowerExport(el []simulator.Evse) float64 {
	var t float64

	for _, e := range el {
		for _, c := range e.Connectors {
			t += c.TPower
		}
	}

	return t
}

/*
Calculate the energy of the CP and retuens tha value (float64)

tp	-	CP total power (float64)

st	-	Simulators start timestamp (time.Time)
*/
func CalculateCPEnergy(tp float64, st time.Time) float64 {
	var t = tp * time.Since(st).Hours()

	return t
}

func IsDataChanClosed(ch <-chan common.Channel) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func IsChanClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func GetStringPointer(s string) *string {
	return &s
}

func GetBoolFromString(v string) bool {
	var b, err = strconv.ParseBool(v)

	if err != nil {
		common.Log("GetBoolFromString").Fatal(err)
	}

	return b
}
