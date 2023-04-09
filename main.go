package main

import (
	"fmt"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
)

func main() {
	var fl = flags.SetFlags()

	fmt.Println(fl)

	common.StartLog()

	_, g := settings.Settings(&fl)

	t := translation.Translation{
		L: translation.Translation{}.GetKey(g.Lang),
	}
	fmt.Println(t.Get(text.MethodNotAllowed))

	var a = api.API{
		Lang: t,
	}

	var srv = a.New()

	a.Serve(srv)
}
