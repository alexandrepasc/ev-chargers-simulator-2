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

	var s = flag.String("s", common.DefSimIP, "Simulator host ip address")

	var t = flag.Int64("t", common.DefTimeout, "Connection timeout (seconds)")

	var csi = flag.String("csi", common.DefCSIP, "Central system ip address")

	var csp = flag.String("csp", common.DefCSPort, "Central system port")

	var gs = flag.String("gs", common.DefGSPath, "General configuration folder, store the application general configurations")

	var ss = flag.String("ss", common.DefSCPath, "Simulators configuration folder, store the simulators configurations")

	var v = flag.Bool("v", false, "")

	flag.Parse()

	if *v {
		fmt.Println("Version:", common.Version)
		os.Exit(0)
	}

	f.HostAddr = *s
	f.ConnTimeout = *t
	f.CSAddr = *csi
	f.CSPort = *csp
	f.GCFolder = *gs
	f.SCFolder = *ss

	return f
}
