package ocpp201

import (
	"sync"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/google/uuid"
	ocpp201 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

type Ocpp201 struct {
	lock       sync.RWMutex                    // Lock goroutine
	logger     logging                         // Logging
	L          translation.Translation         // Translation module
	Timeout    int64                           // Connection timeout
	CSAddr     string                          // Central system ip address
	CSPort     string                          // Central system port
	Asset      *simulator.Asset                // Asset data for the simulator
	Mod        *model.Struct                   // Model data for the asset
	s          ocpp201.ChargingStation         // Ocpp charging station server
	components map[string]component            // Components and variables keys
	bootStatus provisioning.RegistrationStatus // The booting status of the cp
	st         time.Time                       // Simulator start timestamp
	connectSeq bool                            // To trigger the websocket connection
	bootSeq    bool                            // To trigger the boot up sequence
	bootReason provisioning.BootReason         // Boot reason
}

/*
Initialize and setup the requirements to start the server. It will keep the routine running,
handle the channels, and trigger the services and logic to make the simulator work.

c	-	Channel used to communicate the asset information (chan common.Channel)

q	-	Quit channel used to stop the logic and routine(s) (chan bool)
*/
func (o *Ocpp201) Start(c chan common.Channel, q chan bool) {
	o.lock.Lock()

	o.logger = logging{
		toFile: false,
		file:   common.DefGSPath,
	}

	var lm = map[string]string{
		"protocol":  string(o.Asset.Protocol),
		"function":  "Start",
		"simulator": o.Asset.Name,
	}

	o.setStartUpConfigurations()

	// TODO: remove this
	ocppj.SetLogger(logger)
	ws.SetLogger(logger)

	for {
		if o.connectSeq {
			o.s = setupServer(o.Asset.CPId, o.Timeout, o, o.L)

			var conn = "ws://"
			if o.Asset.TLS {
				conn = "wss://"
			}

			sErr := o.s.Start(conn + o.CSAddr + ":" + o.CSPort)

			if sErr != nil {
				o.logger.log(lm, sErr, assets.Error)
				return
			}

			o.connectSeq = false

			o.logger.log(lm, o.L.Get(text.Ocpp201ServerStarted), assets.Info)
		}

		if o.bootSeq {
			var br, eB = o.sendBootNotification(o.bootReason)
			if eB == nil {
				o.processBootResponse(br, o.bootReason)
			}
		}

		select {
		case <-q:
			o.logger.log(lm, o.L.Get(text.Ocpp201ServerStopped), assets.Info)

			o.s.Stop()

			close(q)
			close(c)

			return

		default:
			channelComm(c, o.Asset.Evses, o.st, o.Asset.Name, o.Asset.SimID)
		}
	}
}

/*
Create the websocket and the charge station server, define the configurations for the charge
station and the handler. Returns the server after (ocpp201.ChargingStation).

id	-	Charge station identifier (string)

t	-	Timeout value to set to the server (int64)

o	-	Ocpp201 project structure (*Ocpp201)

l	-	Translation language (translation.Translation)
*/
func setupServer(id string, t int64, o *Ocpp201, l translation.Translation) (s ocpp201.ChargingStation) {
	if o.Asset.TLS {
		s = ocpp201.NewChargingStation(id, nil, assets.GetTLSWsClient(
			t,
			o.Mod.CA,
			o.Mod.Cert,
			o.Mod.Key,
			o.components["SecurityCtrlr"].variables["Identity"].item[0].AttributeValue,
			o.components["SecurityCtrlr"].variables["BasicAuthPassword"].item[0].AttributeValue,
			o.Asset.BasicAuth,
			l,
		))
	} else {
		s = ocpp201.NewChargingStation(id, nil, assets.GetWsClient(
			t,
			o.components["SecurityCtrlr"].variables["Identity"].item[0].AttributeValue,
			o.components["SecurityCtrlr"].variables["BasicAuthPassword"].item[0].AttributeValue,
			o.Asset.BasicAuth,
		))
	}

	s.SetProvisioningHandler(o)
	s.SetRemoteControlHandler(o)

	return s
}

/*
Sends the selected simulator metrics to the channel, to do this the calculations are done on the
fly.

c	-	Communication channel (chan common.Channel)

e	-	Simulator evses list ([]simulator.Evse)

st	-	The start simulator timestamp (time.Time)

n	-	The simulator name (string)

i	-	Simulator identifier (uuid.UUID)
*/
func channelComm(c chan common.Channel, e []simulator.Evse, st time.Time, n string, i uuid.UUID) {
	var tp = assets.CalculateCPPower(e)

	// TODO: the way to find the power factor is not done correctly
	var pf int64

	for _, ei := range e {
		for _, ci := range ei.Connectors {
			if !ci.Enabled {
				continue
			}

			pf = ci.Data[ci.DP.Position].PowerFactor
		}
	}

	c <- common.Channel{
		Name:        n,
		UUID:        i,
		Power:       assets.CalculateInstantCPPower(e),
		PowerFactor: pf,
		Energy:      assets.CalculateCPEnergy(tp, st),
	}
}
