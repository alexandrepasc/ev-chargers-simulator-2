package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
)

/*
Handle the flags supported by the application, with documentation for each one as the default values.

Returns flags structure.
*/
func SetFlags() Flags {
	f := Flags{}

	var s   = common.DefSimIP
	var t   = common.DefTimeout
	var csi = common.DefCSIP
	var csp = common.DefCSPort
	var gs  = common.DefGSPath
	var ss  = common.DefSCPath
	var v bool

	var i = flag.Bool("i", false, "Force the update of the configurations with the values in the flags")

	flag.StringVar(&s, "s", common.DefSimIP, "Simulator host ip address")

	flag.Int64Var(&t, "t", common.DefTimeout, "Connection timeout (seconds)")

	flag.StringVar(&csi, "csi", common.DefCSIP, "Central system ip address")

	flag.StringVar(&csp, "csp", common.DefCSPort, "Central system port")

	flag.StringVar(&gs, "gs", common.DefGSPath, "General configuration folder, store the application general configurations")

	flag.StringVar(&ss, "ss", common.DefSCPath, "Simulators configuration folder, store the simulators configurations")

	flag.BoolVar(&v, "v", false, "Return the current appication version")

	flag.Parse()

	if isFlagPassed("v") {
		fmt.Println("Version:", common.Version)
		os.Exit(0)
	}

	f.ForceUpdate = *i
	
	if isFlagPassed("s") { f.HostAddr = s }
	
	if isFlagPassed("t") {
		f.ConnTimeout = t
	} else {
		f.ConnTimeout = -1
	}

	if isFlagPassed("csi") { f.CSAddr = csi }

	if isFlagPassed("csp") { f.CSPort = csp }

	if isFlagPassed("gs") { f.GCFolder = gs }

	if isFlagPassed("ss") { f.SCFolder = ss }

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
