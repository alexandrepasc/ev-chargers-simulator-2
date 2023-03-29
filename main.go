package main

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
)

func main() {
	var _ = flags.SetFlags()

	common.StartLog()
}
