package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/assets/modbus"
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
	HostIP  string                  // Application host ip address
	Tout    int64                   // Timeout configuration
	Channel []chan common.Channel   // Channel that will enable the communication between the sims
	Info    []assets.DataInfo       // Running assets information
	Quit    []chan bool             // Channel that is used to stop the routines
	stop    chan bool
}

// TODO: add logic to handle the models
/*
Build the channels, the simulators routines, and start the routines.

It returns an array of boolean channels ([]chan bool), one for each simulator running.
*/
func (h *Handler) Start() []chan bool {
	h.Channel = make([]chan common.Channel, len(h.Al))
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

	if getPmIndex(h.Al) != -1 {
		var i = getPmIndex(h.Al)

		var m = getModel(h.Al[i].Model, h.Al[i].Protocol, h.Ml)

		var s = modbus.Modbus{
			L:       h.L,
			HostIP:  h.HostIP,
			Timeout: h.Tout,
			Asset:   h.Al[i],
			Mod:     m,
			Info:    &h.Info,
		}

		go s.Start(h.Channel, h.Quit[i])
	} else {
		h.stop = make(chan bool)
		go h.receiver(h.Channel, h.stop)
	}

	return h.Quit
}

/*
Stops all the simulators routines using the Quit channel array.
*/
func (h *Handler) Stop() {
	for i := range h.Quit {
		h.Quit[i] <- true
	}

	if getPmIndex(h.Al) == -1 {
		h.stop <- true
	}

	for i := range h.Al {
		h.Info[i].UUID = h.Al[i].SimID
		h.Info[i].Name = h.Al[i].Name
		h.Info[i].Status = assets.Inactive
		h.Info[i].Power = 0
		h.Info[i].Energy = 0
	}
}

func (h *Handler) GetStatus() Status {
	var resp Status

	var pmi = getPmIndex(h.Al)

	var a []Assets
	if pmi > -1 {
		a = make([]Assets, len(h.Al)-1)

		resp.Total = int64(len(h.Al) - 1)
	} else {
		a = make([]Assets, len(h.Al))

		resp.Total = int64(len(h.Al))
	}

	for i := range h.Al {
		if pmi == i {
			continue
		}

		var aux = Assets{
			ID:     h.Info[i].UUID,
			Name:   h.Info[i].Name,
			State:  string(h.Info[i].Status),
			Power:  h.Info[i].Power,
			Energy: h.Info[i].Energy,
		}

		a[i-1] = aux
	}

	resp.Assets = a

	return resp
}

/*
Mock a receiver to test the channel communication.
*/
func (h *Handler) receiver(cl []chan common.Channel, s chan bool) {
	var keepOn = true

	for keepOn {
		for i := range cl {
			select {
			case <-s:
				keepOn = false
			case <-cl[i]:
				msg := <-cl[i]
				h.Info[i] = assets.DataInfo{
					UUID:        msg.UUID,
					Name:        msg.Name,
					Status:      assets.Active,
					Power:       msg.Power,
					PowerFactor: msg.PowerFactor,
					Energy:      msg.Energy,
				}

				fmt.Println(msg)
			default:
				continue
			}
		}

		time.Sleep(1 * time.Second)
	}
}

/*
Initialize the info variable and set the starting values from all the assets loaded to the
variable.
*/
func (h *Handler) SetInfoStartValues() {
	h.Info = make([]assets.DataInfo, len(h.Al))

	for ai, a := range h.Al {
		h.Info[ai].UUID = a.SimID
		h.Info[ai].Name = a.Name
		h.Info[ai].Status = assets.Inactive
		h.Info[ai].Power = 0
		h.Info[ai].PowerFactor = 0
		h.Info[ai].Energy = 0
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

func getPmIndex(al []*simulator.Asset) int {
	for i, a := range al {
		if a.Type == simulator.Pm {
			return i
		}
	}

	return -1
}
