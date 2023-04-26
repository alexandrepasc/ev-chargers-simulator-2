package handler

import (
	"fmt"
	"strings"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets/ocpp16"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/google/uuid"
)

const buf = 10

type Handler struct {
	L       translation.Translation // Translation language settings
	Al      []*simulator.Asset      // Assets list
	Ml      []*model.Struct         // Ocpp models list
	Addr    string                  // Central system ip address
	Port    string                  // Central system port
	Tout    int64                   // Timeout configuration
	Channel chan common.Channel     // Channel that will enable the communication between the sims
	Quit    []chan bool             // Channel that is used to stop the routines
}

// TODO: add logic to handle the models
/*
Build the channels, the simulators routines, and start the routines.

It returns an array of boolean channels ([]chan bool), one for each simulator running.
*/
func (h *Handler) Start() []chan bool {
	h.Channel = make(chan common.Channel, len(h.Al)+buf)
	h.Quit = make([]chan bool, len(h.Al))

	for i, a := range h.Al {
		h.Quit[i] = make(chan bool, 1)

		var m = getModel(a.Model, a.Protocol, h.Ml)

		switch a.Protocol {
		case simulator.Ocpp16:
			var s = ocpp16.Ocpp16{
				L:       h.L,
				Timeout: h.Tout,
				CSAddr:  h.Addr,
				CSPort:  h.Port,
				Asset:   a,
				Mod:     m,
			}

			go s.Start(h.Channel, h.Quit[i])
		case simulator.Ocpp201:

		case simulator.Modbus:
		}
	}

	go receiverName(h.Channel)

	return h.Quit
}

/*
Stops all the simulators routines using the Quit channel array.
*/
func (h *Handler) Stop() {
	for i := range h.Al {
		h.Quit[i] <- true
	}
}

/*
Mock a simulator routine to test the channel communication.
*/
// func runSimulator(s sim, c chan common.Channel, q chan bool) {
// 	for i := 0; i < 20; i++ {
// 		time.Sleep(1 * time.Second)

// 		select {
// 		case <-q:
// 			return
// 		default:
// 			c <- common.Channel{
// 				Name: s.asset.Name,
// 				Uuid: s.asset.SimID,
// 			}
// 		}
// 	}
// }

/*
Mock a receiver to test the channel communication.
*/
func receiverName(c chan common.Channel) {
	for msg := range c {
		fmt.Println(msg)
	}
}

// TODO: Handle modbus protocol
/**/
func getModel(id uuid.UUID, p simulator.Protocol, al []*model.Struct) *model.Struct {
	for _, m := range al {
		if m.ID == id {
			if strings.Contains(string(p), string(m.Type)) {
				return m
			}
		}
	}

	switch p {
	case simulator.Ocpp16:
		return &model.DefOcppMod
	case simulator.Ocpp201:
		return &model.DefOcppMod
	case simulator.Modbus:
		return nil
	default:
		return nil
	}
}
