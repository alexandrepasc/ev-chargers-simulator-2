package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
)

/*
	Handle the flags supported by the application, with documentation for each one, as the default values
*/
func SetFlags() (flags) {

	f := flags{}

	var s *string = flag.String("s", "84.5.7.64", "Simulator host ip address")
	var t *int = flag.Int("t", 70, "Connection timeout (seconds)")
	var csi *string = flag.String("csi", "iot-gate-imx8.lan", "Central system ip address")
	var csp *string = flag.String("csp", "49443", "Central system port")
	var c *string = flag.String("c", "", "Configuration folder, will ignore files with \"_\" as 1st char")
	var v *bool = flag.Bool("v", false, "")

	flag.Parse()

	if *v {
		fmt.Println("Version:", common.Version)
		os.Exit(0)
	}

	f.HostAddr = *s
	f.ConnTimeout = *t
	f.CSAddr = *csi
	f.CSPort = *csp
	f.ConfFolder = *c

	return f
}