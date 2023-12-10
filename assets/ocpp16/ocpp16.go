package ocpp16

import (
	"sync"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/google/uuid"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type Ocpp16 struct {
	lock      sync.RWMutex                     // Lock goroutine
	logger    logging                          // Logging
	L         translation.Translation          // translation
	Timeout   int64                            // Connection timeout
	CSAddr    string                           // Central system ip address
	CSPort    string                           // Central system port
	Asset     *simulator.Asset                 // Asset data for the simulator
	Mod       *model.Struct                    // Model data for the asset
	s         ocpp16.ChargePoint               // Ocpp charge point server
	Conf      map[string]core.ConfigurationKey // Configuration key map
	bootSeq   assets.BootSeq                   // Boot sequence structure
	resetSeq  assets.ResetSeq                  // Reset sequence structure
	localAuth struct {                         // Local auth list
		version int64                         // Version identifier
		list    []localauth.AuthorizationData // List with the authorization information
	}
	chargeProfile  *types.ChargingProfile // Charging profile set by the CS
	txnAlignedData []types.MeterValue     // Store the transaction aligned data
	txnSampledData []types.MeterValue     // store the transaction sampled data
	tick           int64                  // Ticker to enable trigger scheduled events
	heartbeatC     int64                  // Heartbeat count to handle the request interval
	st             time.Time
}

/**/
func (o *Ocpp16) Start(c chan common.Channel, q chan bool) {
	o.lock.Lock()

	o.logger = logging{
		toFile: false,
		file:   common.DefGSPath,
	}

	o.setStartUpConfigurations()

	o.s = setupServer(o.Asset.CPId, o.Timeout, o, o.L)

	sErr := o.s.Start(assets.GetConnProtocol(o.Asset.TLS) + o.CSAddr + ":" + o.CSPort)

	if sErr != nil {
		o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start"}, sErr.Error(), assets.Error)
		return
	}

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start", "simulator": o.Asset.Name}, "Ocpp 1.6 server started", assets.Info)

	for {
		if o.bootSeq.IsToTrigger {
			if o.bootSeq.BootInterval == 0 || o.tick%int64(o.bootSeq.BootInterval) == 0 {
				var br, be = o.sendBootNotification()
				if be == nil {
					o.processBootResponse(br)
				}
			}
		}

		go o.updateData()

		go o.processSampledData()

		o.processAlignedData()

		o.processHeartbeat()

		o.tick = assets.HandleTick(o.tick)

		select {
		case <-q:
			o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start", "simulator": o.Asset.Name}, "Ocpp 1.6 server stop", assets.Info)

			o.s.Stop()

			close(q)
			close(c)

			return
		default:
			channelComm(c, o.Asset.Evses, o.st, o.Asset.Name, o.Asset.SimID)
		}

		time.Sleep(1 * time.Second)
	}
}

/*
Create the websocket and the charge station server, define the configurations for the charge
station and the handler. Returns the server after (ocpp16.ChargingStation).

id	-	Charge station identifier (string)

t	-	Timeout value to set to the server (int64)

o	-	Ocpp16 project structure (*Ocpp16)

l	-	Translation language (translation.Translation)
*/
func setupServer(id string, t int64, o *Ocpp16, l translation.Translation) (s ocpp16.ChargePoint) {
	// The basic auth is using the username and password directly from the Model
	if o.Asset.TLS {
		s = ocpp16.NewChargePoint(id, nil, assets.GetTLSWsClient(
			t,
			o.Mod.CA,
			o.Mod.Cert,
			o.Mod.Key,
			o.Mod.BasicAuth.Username,
			o.Mod.BasicAuth.Password,
			o.Asset.BasicAuth,
			l,
		))
	} else {
		s = ocpp16.NewChargePoint(id, nil, assets.GetWsClient(
			t,
			&o.Mod.BasicAuth.Username,
			&o.Mod.BasicAuth.Password,
			o.Asset.BasicAuth,
		))
	}

	s.SetCoreHandler(o)
	s.SetLocalAuthListHandler(o)
	s.SetRemoteTriggerHandler(o)

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
// TODO: the calculate cp power needs a list of evses and it can be useful to the ocpp 2.0.1
func channelComm(c chan common.Channel, e []simulator.Evse, st time.Time, n string, i uuid.UUID) {
	var tp = assets.CalculateCPPower(e)

	// TODO: the way to find the power factor is not done correctly
	var pf int64

	for _, ci := range e[0].Connectors {
		if !ci.Enabled {
			continue
		}

		pf = ci.Data[ci.DP.Position].PowerFactor
	}

	c <- common.Channel{
		Name:        n,
		UUID:        i,
		Power:       assets.CalculateInstantCPPower(e),
		PowerFactor: pf,
		Energy:      assets.CalculateCPEnergy(tp, st),
	}
}
