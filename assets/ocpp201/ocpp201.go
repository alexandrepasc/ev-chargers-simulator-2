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
)

type Ocpp201 struct {
	lock   sync.RWMutex            // Lock goroutine
	logger logging                 // Logging
	L      translation.Translation // translation
	CSAddr string                  // Central system ip address
	CSPort string                  // Central system port
	Asset  *simulator.Asset        // Asset data for the simulator
	Mod    *model.Struct           // Model data for the asset
	st     time.Time               // Simulator start timestamp
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

	for {
		select {
		case <-q:
			o.logger.log(
				map[string]string{
					"protocol":  string(o.Asset.Protocol),
					"function":  "Start",
					"simulator": o.Asset.Name,
				},
				o.L.Get(text.Ocpp201ServerStopped), assets.Info)

			// o.s.Stop()

			close(q)
			close(c)

			return

		default:
			channelComm(c, o.Asset.Evses, o.st, o.Asset.Name, o.Asset.SimID)
		}
	}
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
