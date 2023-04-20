package handler

import (
	"fmt"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
)

const buf = 10

type Handler struct {
	L       translation.Translation // Translation language settings
	Al      []*simulator.Asset      // Assets list
	Oml     []model.OcppModel       // Ocpp models list
	Sims    []sim                   // List of simulators
	Channel chan channel            // Channel that will enable the communication between the sims
	Quit    []chan bool             // Channel that is used to stop the routines
}

// TODO: add logic to handle the models
/*
Build the channels, the simulators routines, and start the routines.

It returns an array of boolean channels ([]chan bool), one for each simulator running.
*/
func (h *Handler) Start() []chan bool {
	h.Channel = make(chan channel, len(h.Al)+buf)
	h.Quit = make([]chan bool, len(h.Al))

	for i, a := range h.Al {
		var sa = sim{
			asset: a,
		}

		h.Quit[i] = make(chan bool, 1)

		h.Sims = append(h.Sims, sa)
	}

	for i, s := range h.Sims {
		go runSimulator(s, h.Channel, h.Quit[i])
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
func runSimulator(s sim, c chan channel, q chan bool) {
	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)

		select {
		case <-q:
			return
		default:
			c <- channel{
				name: s.asset.Name,
				uuid: s.asset.SimID,
			}
		}
	}
}

/*
Mock a receiver to test the channel communication.
*/
func receiverName(c chan channel) {
	for msg := range c {
		fmt.Println(msg)
	}
}
