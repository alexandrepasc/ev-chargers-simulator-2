package ocpp16

import (
	"sync"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ws"
)

type Ocpp16 struct {
	lock    sync.RWMutex            // Lock goroutine
	logger  logging                 // Logging
	L       translation.Translation // translation
	Timeout int64                   // Connection timeout
	CSAddr  string                  // Central system ip address
	CSPort  string                  // Central system port
	Asset   *simulator.Asset        // Asset data for the simulator
	Mod     *model.Struct           // Model data for the asset
	s       ocpp16.ChargePoint
}

/**/
func (o *Ocpp16) Start(c chan common.Channel, q chan bool) { //nolint:revive // because dev
	o.lock.Lock()

	o.logger = logging{
		toFile: false,
		file:   common.DefGSPath,
	}

	var h = &Handler{}

	o.s = setupServer(o.Asset.CPId, o.Timeout, h)

	sErr := o.s.Start("ws://" + o.CSAddr + ":" + o.CSPort)

	if sErr != nil {
		o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start"}, sErr.Error(), assets.Error)
		return
	}

	o.logger.log(map[string]string{"protocol": "ocpp1.6", "function": "Start", "simulator": o.Asset.Name}, "Ocpp 1.6 server started", assets.Info)

	o.sendBootNotification()
}

/**/
func setupServer(id string, t int64, h *Handler) (s ocpp16.ChargePoint) {
	s = ocpp16.NewChargePoint(id, nil, getWsClient(t))

	s.SetCoreHandler(h)

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
