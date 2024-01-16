package modbus

import (
	"fmt"
	"sync"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/simonvetter/modbus"
	"github.com/sirupsen/logrus"
)

type Modbus struct {
	lock    sync.RWMutex            // Lock goroutine
	logger  assets.Logging          // Logging
	L       translation.Translation // translation
	HostIP  string                  // Application host ip address
	Timeout int64                   // Connection timeout
	Asset   *simulator.Asset        // Asset data for the simulator
	Mod     *model.Struct           // Model data for the asset
	Info    *[]assets.DataInfo      // Model with the assets information
	s       *modbus.ModbusServer    // Modbus server
}

/*
The start and stop function for the modbus server. It is here that all the logic will be directly
or indirectly be called.

c	-	List of channels from the other assets in PM case, or single channel in EVC case ([]chan common.Channel)

q	-	Quit channel to be able to stop the routine from the main application (chan bool)
*/
func (m *Modbus) Start(c []chan common.Channel, q chan bool) {
	m.lock.Lock()

	m.logger = assets.Logging{
		ToFile: false,
		File:   common.DefGSPath,
		Logger: logrus.New(),
		L:      m.L,
	}

	if m.Asset.LogLevel != nil {
		m.logger.Level = assets.TranslateLogLevels(*m.Asset.LogLevel)
	} else {
		m.logger.Level = logrus.ErrorLevel
	}

	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "Start",
		"simulator": m.Asset.Name,
	}

	m.s = setupServer(m)

	var err = m.s.Start()
	if err != nil {
		m.logger.Log(lm, err, assets.Error)
		return
	}

	m.logger.Log(lm, m.L.Get(text.ModbusServerStarted), assets.Info)

	for {
		select {
		case <-q:
			var errS = m.s.Stop()

			if errS != nil {
				m.logger.Log(lm, errS, assets.Fatal)
			}

			m.logger.Log(lm, m.L.Get(text.ModbusServerStopped), assets.Info)

			close(q)

			return
		default:
			break
		}

		switch m.Asset.Type {
		case simulator.Pm:
			m.pmCycle(c)

		case simulator.Evc:
		}

		time.Sleep(1 * time.Second)
	}
}

/*
Have all the logic dedicated to the PM asset type. This will be called when the asset has the
PM type and will have all the logic for this asset type.

c	-	List of channels from the other assets ([]chan common.Channel)
*/
func (m *Modbus) pmCycle(c []chan common.Channel) {
	for i := range c {
		select {
		case <-c[i]:
			var msg = <-c[i]

			var aux = assets.DataInfo{
				UUID:        msg.UUID,
				Name:        msg.Name,
				Status:      assets.Active,
				Power:       msg.Power,
				PowerFactor: msg.PowerFactor,
				Energy:      msg.Energy,
			}

			(*m.Info)[i] = aux

			fmt.Println(msg)
		default:
			continue
		}
	}
}

/*
Setup the modbus server configurations, will receive the needed variables and will return the
modbus server structure (*modbus.ModbusServer). In case it fails will log a fatal error to the
console.

m	-	The current file structure (*Modbus)
*/
func setupServer(m *Modbus) *modbus.ModbusServer {
	var serv, err = modbus.NewServer(
		&modbus.ServerConfiguration{
			URL:     "tcp://" + m.HostIP + ":" + m.Asset.Port,
			Timeout: time.Second * time.Duration(m.Timeout),
		},
		m,
	)

	if err != nil {
		var lm = map[string]string{
			"protocol":  string(m.Asset.Protocol),
			"function":  "setupServer",
			"simulator": m.Asset.Name,
		}

		m.logger.Log(lm, err, assets.Fatal)
	}

	return serv
}
