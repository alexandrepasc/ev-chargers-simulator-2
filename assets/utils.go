package assets

import (
	"crypto/tls"
	"crypto/x509"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/sirupsen/logrus"
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
Call the calculate current function for all the evses and all the connectors, sums the values, and
returns that value (float64).
*/
func CalculateCPTotalCurrent(evses []simulator.Evse, phn int64) (t float64) {
	for ei, e := range evses {
		for ci := range e.Connectors {
			var v int64

			for _, vi := range evses[ei].Connectors[ci].Data[evses[ei].Connectors[ci].DP.Position].Voltage {
				if vi == 0 {
					continue
				}

				v = vi
			}

			t += CalculateCurrent(
				evses[ei].Connectors[ci].Data[evses[ei].Connectors[ci].DP.Position].Power,
				evses[ei].Connectors[ci].Data[evses[ei].Connectors[ci].DP.Position].PowerFactor,
				v,
				phn,
			)
		}
	}

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
Sum the power of all the CP EVSEs active connectors and returns the total value (float64).

el	-	The list of evses that the CP has ([]simulator.Evse)
*/
func CalculateInstantCPPower(el []simulator.Evse) float64 {
	var t float64

	for ei, e := range el {
		for ci, c := range e.Connectors {
			if !c.Enabled {
				continue
			}

			t += float64(el[ei].Connectors[ci].Data[el[ei].Connectors[ci].DP.Position].Power)
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

/*
Calculate the max current sent to the ev by the cp. Will loop all the connector data,
calculating the current for each of the data, and comparing if it is the max current. At
the end it will return the max current value (float64)

c	-	Connector (simulator.Connector)

ph	-	CP number of phases (int)

vi	-	Phase index (int)
*/
func CalculateConnectorMaxCurrent(c *simulator.Connector, ph, vi int) float64 {
	var mc float64

	for _, i := range c.Data {
		var a = CalculateCurrent(i.Power, i.PowerFactor, i.Voltage[vi], int64(ph))

		if a > mc {
			mc = a
		}
	}

	return mc
}

/*
Calculate the max current sent to the ev by the cp. Will loop all the connectors if a connector
is enabled will loop through the data. For each data will canculate the current and compare it with
the previous stored current. At the end will return the max current from the data list (float64)
or 0 in case no connector is enabled.

e	-	Evse structure (simulator.Evse)

ph	-	Asset phase number (int)
*/
func CalculateMaxCurrent(e *simulator.Evse, ph int) float64 {
	var mc float64

	for _, c := range e.Connectors {
		if !c.Enabled {
			continue
		}

		for _, i := range c.Data {
			var v int64

			for _, vi := range i.Voltage {
				if vi > 0 {
					v = vi

					break
				}
			}

			var a = CalculateCurrent(i.Power, i.PowerFactor, v, int64(ph))

			if a > mc {
				mc = a
			}
		}

		return mc
	}

	return 0
}

/*
Get the max power sent to the ev by the cp connector. Will loop all the connector data and
comparing if it with the previous stored one. Will return the max power from the data list (int64).

c	-	Connector structure (*simulator.Connector)
*/
func GetConnectorMaxPower(c *simulator.Connector) int64 {
	var p int64

	for _, d := range c.Data {
		if d.Power < p {
			continue
		}

		p = d.Power
	}

	return p
}

/*
Get the max power sent to the ev by the cp. Will loop all the connectors if the connector is
enabled loop through it's data and compare it with the previous stored one. Will return the
max power from the data list (int64).

e	-	Evse structure (*simulator.Evse)
*/
func GetMaxPower(e *simulator.Evse) int64 {
	var p int64

	for _, c := range e.Connectors {
		if !c.Enabled {
			continue
		}

		for _, d := range c.Data {
			if d.Power < p {
				continue
			}

			p = d.Power
		}
	}

	return p
}

/*
Create a new websocket client with the parameter timeout and returns it (*ws.Client).

t	-	Timeout value (int64)

u		-	Basic authentication username (string)

p		-	Basic authentication password (string)

ba	-	If the client should have the http basic authentication (bool)
*/
func GetWsClient(t int64, u, p *string, ba bool) (wsc *ws.Client) {
	wsc = ws.NewClient()

	if ba {
		wsc.SetBasicAuth(*u, *p)
	}

	var cfg = ws.ClientTimeoutConfig{
		HandshakeTimeout: time.Second * time.Duration(t),
		WriteWait:        time.Second * time.Duration(t),
		PingPeriod:       time.Second * time.Duration(t),
		PongWait:         time.Second * time.Duration(t),
	}

	wsc.SetTimeoutConfig(cfg)

	return wsc
}

/*
Create a new tls websocket client with the parameters sent and returns it (*ws.Client).

t		-	Timeout value (int64)

ca		-	CA certificate path and name (string)

cert	-	Client certificate path and name (string)

key		-	Client certificate key path and name (string)

u		-	Basic authentication username (string)

p		-	Basic authentication password (string)

ba		-	If the client should have the http basic authentication (bool)

l		-	Translation language (translation.Translation)
*/
func GetTLSWsClient(t int64, ca, cert, key, u, p string, ba bool, l translation.Translation) (wsc *ws.Client) {
	var certPool, errcp = x509.SystemCertPool()
	if errcp != nil {
		common.Log("GetTlsWsClient").Error(errcp)
	}

	var caCert, errca = os.ReadFile(ca)
	if errca != nil {
		common.Log("GetTlsWsClient").Fatal(errca)
	} else if !certPool.AppendCertsFromPEM(caCert) {
		common.Log("GetTlsWsClient").Info(l.Get(text.CaCertNotFound))
	}

	var clientCertificates []tls.Certificate

	var certificate, errc = tls.LoadX509KeyPair(cert, key)
	if errc != nil {
		common.Log("GetTlsWsClient").Fatal(errc)
	}

	clientCertificates = []tls.Certificate{certificate}

	wsc = ws.NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: clientCertificates,
		MinVersion:   tls.VersionTLS13,
	})

	if ba {
		wsc.SetBasicAuth(u, p)
	}

	var cfg = ws.ClientTimeoutConfig{
		HandshakeTimeout: time.Second * time.Duration(t),
		WriteWait:        time.Second * time.Duration(t),
		PingPeriod:       time.Second * time.Duration(t),
		PongWait:         time.Second * time.Duration(t),
	}

	wsc.SetTimeoutConfig(cfg)

	return wsc
}

/*
Handles the counter from the simulator and handles the max int64 value, in case it is reaching the
max value (max int64 value - 7) it will be reseted to 0 and will return the value (int64).

t	-	Ticker value (int64)
*/
func HandleTick(t int64) int64 {
	const rInt64 = math.MaxInt64 - 7

	if t > rInt64 {
		t = 0
	} else {
		t++
	}

	return t
}

/*
Get the connection protocol string depending if the tls is enabled for the asset or not (string).

t	-	Is the tls activated to the asset (bool)
*/
func GetConnProtocol(t bool) string {
	if t {
		return string(wssConn)
	}

	return string(wsConn)
}

/*
Runs through all the asset EVSEs and retursn true (bool) in case one of the connectors is in the
charging state, if not returns false.

el 	-	Asset EVSEs list ([]*simulator.Evse)
*/
func IsAssetWithTransaction(el []simulator.Evse) bool {
	for _, e := range el {
		if IsEvseWithTransaction(e) {
			return true
		}
	}

	return false
}

/*
Returns true (bool) in case any of the EVSE connector is in the charging state and false if
none of the connectors is in the that state.

e	-	Asset EVSE (*simulator.Evse)
*/
func IsEvseWithTransaction(e simulator.Evse) bool {
	for _, c := range e.Connectors {
		if c.Data[c.DP.Position].ChargingState == int64(Charging) {
			return true
		}
	}

	return false
}

/*
Get the EVSE index and structure (int, *simulator.Evse) from the asset using the identifier. In
case none is found or the id is nil returns the index as -1 and an empty EVSE structure.

id	-	The EVSE identifier (*int)

el	-	The asset EVSE list ([]*simulator.Evse)
*/
func GetEvseByID(id *int, el []simulator.Evse) (i int, evse simulator.Evse) {
	if id == nil {
		return -1, simulator.Evse{}
	}

	for idx, e := range el {
		if e.ID == int64(*id) {
			return idx, e
		}
	}

	return -1, simulator.Evse{}
}

/*
Gets the index (*int) and the structure (simulator.Evse) of the active evse. If there is no
active evse returns nil for the index and an empty structure.

el	-	The asset EVSE list ([]simulator.Evse)
*/
func GetActiveEvse(el []simulator.Evse) (idx *int, evse simulator.Evse) {
	for i, ie := range el {
		if IsEvseWithTransaction(ie) {
			return &i, ie
		}
	}

	return nil, simulator.Evse{}
}

/*
Get the active connector of given EVSE from an asset. It will return the connector index (*int)
and the connector structure (*simulator.Connector) of the active one. If there is no active
connector it returns nil in both.

e	-	Asset evse structure (*simulator.Evse)
*/
func GetActiveConnector(e *simulator.Evse) (idx *int, c *simulator.Connector) {
	for i, ic := range e.Connectors {
		if e.Connectors[i].Data[e.Connectors[i].DP.Position].ChargingState == int64(Charging) {
			return &i, &ic
		}
	}

	return nil, nil
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

/*
Parse a string value to a boolean variable and returns it (bool).

Triggers a fatal error in case it wasn't able to parse it.

f	-	A string with the function name that called the function (string)

v	-	The string value to be parsed (string)
*/
func GetBoolFromString(f, v string) bool {
	var b, err = strconv.ParseBool(v)

	if err != nil {
		common.Log("GetBoolFromString - " + f).Fatal(err)
	}

	return b
}

/*
Receives a string with from the asset configuration and returns the matching logging level
(logrus.Level). If the string does not match any of the cases it will return the error level.

l	-	Configuration log level information (string)
*/
func TranslateLogLevels(l string) logrus.Level {
	switch l {
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return logrus.ErrorLevel
	}
}
