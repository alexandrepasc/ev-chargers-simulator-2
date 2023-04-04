package settings

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/configs"
)

func Settings(fl *flags.Flags) {
	var _ = configs.Configs(fl)
}
