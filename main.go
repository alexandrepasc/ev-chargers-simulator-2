package main

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/api"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
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

	var a = api.API{
		Lang:    t,
		General: g,
		Sim:     s,
	}

	var srv = a.New()

	a.Serve(srv)
}
