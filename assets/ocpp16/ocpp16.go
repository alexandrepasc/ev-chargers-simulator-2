package ocpp16

import (
	"strconv"
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
	"github.com/lorenzodonini/ocpp-go/ws"
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
	localAuth struct {                         // Local auth list
		version int64                         // Version identifier
		list    []localauth.AuthorizationData // List with the authorization information
	}
	chargeProfile  *types.ChargingProfile // Charging profile set by the CS
	txnAlignedData []types.MeterValue     // Store the transaction aligned data
	txnSampledData []types.MeterValue     // store the transaction sampled data
	t              int64
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

	o.s = setupServer(o.Asset.CPId, o.Timeout, o)

	sErr := o.s.Start("ws://" + o.CSAddr + ":" + o.CSPort)

	if sErr != nil {
		o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start"}, sErr.Error(), assets.Error)
		return
	}

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start", "simulator": o.Asset.Name}, "Ocpp 1.6 server started", assets.Info)

	go o.sendBootNotification()

	for {
		select {
		case <-q:
			o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start", "simulator": o.Asset.Name}, "Ocpp 1.6 server stop", assets.Info)

			o.s.Stop()

			return
		default:
			break
		}

		o.updateData()

		go o.processSampledData()

		var hbi, _ = strconv.ParseInt(*o.Conf["HeartbeatInterval"].Value, 10, 64)
		if o.t%hbi == 0 {
			go o.heartbeat()
		}

		channelComm(c, o.Asset.Evses, o.st, o.Asset.Name, o.Asset.SimID)

		o.processAlignedData()

		o.handleTick()

		time.Sleep(1 * time.Second)
	}
}

/**/
func setupServer(id string, t int64, h *Ocpp16) (s ocpp16.ChargePoint) {
	s = ocpp16.NewChargePoint(id, nil, getWsClient(t))

	s.SetCoreHandler(h)
	s.SetLocalAuthListHandler(h)
	s.SetRemoteTriggerHandler(h)

	return s
}

/**/
func getWsClient(t int64) (wsc *ws.Client) {
	wsc = ws.NewClient()

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

	c <- common.Channel{
		Name:   n,
		UUID:   i,
		Power:  tp,
		Energy: assets.CalculateCPEnergy(tp, st),
	}
}
