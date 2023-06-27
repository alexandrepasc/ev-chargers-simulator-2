package main

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/api"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/handler"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
)

func main() {
	var fl = flags.SetFlags()

	common.StartLog()

	c, g := settings.Settings(&fl)

	t := translation.Translation{
		L: translation.Translation{}.GetKey(g.Lang),
	}

	s := simulator.Simulator{
		Scp: c.SimulatorsConfigFolder,
		L:   t,
	}

	m := model.Model{
		Scp: c.SimulatorsConfigFolder,
		L:   t,
	}

	h := handler.Handler{
		L:    t,
		Addr: g.CSAddr,
		Port: g.CSPort,
		Tout: g.ConnTimeout,
	}

	var a = api.API{
		Lang:        t,
		General:     g,
		GeneralPath: c.GeneralConfigFolder,
		Sim:         s,
		Mod:         m,
		Hand:        h,
	}

	var srv = a.New()

	a.Serve(srv)
}
