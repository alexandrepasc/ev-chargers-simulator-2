package flags

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
)

/*
Handle the flags supported by the application, with documentation for each one as the default values.

Returns flags structure.
*/
func SetFlags() Flags {
	f := Flags{}

	var s = common.DefSimIP

	var t = common.DefTimeout

	var csi = common.DefCSIP

	var csp = common.DefCSPort

	var gs = common.DefGSPath

	var ss = common.DefSCPath

	var l = common.DefLanguage

	var ai = common.DefAPIAddr

	var ap = common.DefAPIPort

	var v bool

	var i = flag.Bool("i", false, "Force the update of the configurations with the values in the flags")

	flag.StringVar(&s, "s", common.DefSimIP, "Simulator host ip address")

	flag.Int64Var(&t, "t", common.DefTimeout, "Connection timeout (seconds)")

	flag.StringVar(&csi, "csi", common.DefCSIP, "Central system ip address")

	flag.StringVar(&csp, "csp", common.DefCSPort, "Central system port")

	flag.StringVar(&gs, "gs", common.DefGSPath, "General configuration folder, store the application general configurations")

	flag.StringVar(&ss, "ss", common.DefSCPath, "Simulators configuration folder, store the simulators configurations")

	flag.StringVar(&l, "l", common.DefLanguage, "Language used by the application ["+strings.Join(translation.Translation{}.List(), ", ")+"]")

	flag.StringVar(&ai, "ai", common.DefAPIAddr, "Address used to serve the api http server")

	flag.StringVar(&ap, "ap", common.DefAPIPort, "")

	flag.BoolVar(&v, "v", false, "Return the current application version")

	flag.Parse()

	if isFlagPassed("v") {
		fmt.Println("Version:", common.Version)
		os.Exit(0)
	}

	f.ForceUpdate = *i

	if isFlagPassed("s") {
		f.HostAddr = s
	}

	if isFlagPassed("t") {
		f.ConnTimeout = t
	} else {
		f.ConnTimeout = -1
	}

	if isFlagPassed("csi") {
		f.CSAddr = csi
	}

	if isFlagPassed("csp") {
		f.CSPort = csp
	}

	if isFlagPassed("gs") {
		f.GCFolder = gs
	}

	if isFlagPassed("ss") {
		f.SCFolder = ss
	}

	if isFlagPassed("l") {
		f.Language = l
	}

	if isFlagPassed("ai") {
		f.APIAddr = ai
	}

	if isFlagPassed("ap") {
		f.APIPort = ap
	}

	return f
}

func isFlagPassed(n string) bool {
	var found = false

	flag.Visit(func(fl *flag.Flag) {
		if fl.Name == n {
			found = true
		}
	})

	return found
}
