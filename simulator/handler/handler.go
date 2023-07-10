package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets/ocpp16"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/google/uuid"
)

type Handler struct {
	L       translation.Translation // Translation language settings
	Al      []*simulator.Asset      // Assets list
	Ml      []*model.Struct         // Ocpp models list
	Addr    string                  // Central system ip address
	Port    string                  // Central system port
	Tout    int64                   // Timeout configuration
	Channel []chan common.Channel   // Channel that will enable the communication between the sims
	Info    []info                  // Running assets information
	Quit    []chan bool             // Channel that is used to stop the routines
}

// TODO: add logic to handle the models
/*
Build the channels, the simulators routines, and start the routines.

It returns an array of boolean channels ([]chan bool), one for each simulator running.
*/
func (h *Handler) Start() []chan bool {
	h.Channel = make([]chan common.Channel, len(h.Al))
	h.Info = make([]info, len(h.Al))
	h.Quit = make([]chan bool, len(h.Al))

	for i, a := range h.Al {
		h.Channel[i] = make(chan common.Channel, cBuf)
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

			for x := range s.Asset.Evses {
				for y := range s.Asset.Evses[x].Connectors {
					s.Asset.Evses[x].Connectors[y].Enabled = false
				}
			}

			go s.Start(h.Channel[i], h.Quit[i])
		case simulator.Ocpp201:

		case simulator.Modbus:
		}
	}

	go h.receiver(h.Channel, h.Quit)

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

func (h *Handler) GetStatus() Status {
	var resp Status

	var a = make([]Assets, len(h.Al))

	resp.Total = int64(len(h.Al))

	for i := range h.Al {
		var aux = Assets{
			ID:     h.Info[i].UUID,
			Name:   h.Info[i].Name,
			State:  string(h.Info[i].Status),
			Power:  h.Info[i].Power,
			Energy: h.Info[i].Energy,
		}

		a[i] = aux
	}

	resp.Assets = a

	return resp
}

/*
Mock a receiver to test the channel communication.
*/
func (h *Handler) receiver(cl []chan common.Channel, ql []chan bool) {
	var keepOn = true

	for keepOn {
		for i := range cl {
			select {
			case <-ql[i]:
				keepOn = false
			case <-cl[i]:
				msg := <-cl[i]
				h.Info[i] = info{
					UUID:   msg.UUID,
					Name:   msg.Name,
					Status: active,
					Power:  msg.Power,
					Energy: msg.Energy,
				}

				fmt.Println(msg)
			default:
				continue
			}
		}

		time.Sleep(1 * time.Second)
	}

	for i := range cl {
		close(ql[i])
		close(cl[i])

		h.Info[i].Status = inactive
		h.Info[i].Power = 0
		h.Info[i].Energy = 0
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
