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
)

type Modbus struct {
	lock    sync.RWMutex            // Lock goroutine
	logger  logging                 // Logging
	L       translation.Translation // translation
	Timeout int64                   // Connection timeout
	Asset   *simulator.Asset        // Asset data for the simulator
	Mod     *model.Struct           // Model data for the asset
	Info    *[]assets.DataInfo
}

/**/
func (m *Modbus) Start(c []chan common.Channel, q chan bool) {
	m.lock.Lock()

	m.logger = logging{
		toFile: false,
		file:   common.DefGSPath,
	}

	var lm = map[string]string{
		"protocol":  string(m.Asset.Protocol),
		"function":  "Start",
		"simulator": m.Asset.Name,
	}

	m.logger.log(lm, m.L.Get(text.ModbusServerStarted), assets.Info)

	for {
		switch m.Asset.Type {
		case simulator.Pm:
			m.pmCycle(c)
		case simulator.Evc:
		}

		select {
		case <-q:
			m.logger.log(lm, m.L.Get(text.ModbusServerStopped), assets.Info)

			close(q)

			return
		default:
			break
		}

		time.Sleep(1 * time.Second)
	}
}

/**/
func (m *Modbus) pmCycle(c []chan common.Channel) {
	for i := range c {
		select {
		case <-c[i]:
			var msg = <-c[i]

			var aux = assets.DataInfo{
				UUID:   msg.UUID,
				Name:   msg.Name,
				Status: assets.Active,
				Power:  msg.Power,
				Energy: msg.Energy,
			}

			(*m.Info)[i] = aux

			fmt.Println(msg)
		default:
			continue
		}
	}
}
