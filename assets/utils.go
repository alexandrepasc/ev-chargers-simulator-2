package assets

import (
	"math"
	"time"
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
func CalculateEnergy(tp float64, st time.Time) (e float64) {
	e = tp * time.Since(st).Hours()

	return e
}
