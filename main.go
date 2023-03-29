package main

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/logging"
)

func main() {

	var _ = flags.SetFlags()

	logging.StartLog()
}