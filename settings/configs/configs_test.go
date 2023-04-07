//nolint:gocritic,nolintlint
package configs_test

import (
	"os"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/configs"
	"github.com/stretchr/testify/assert"
)

const (
	file            string = "/configs.json"
	generalFolder   string = "/general"
	simulatorFolder string = "/simulator"
)

func TestConfigsNoFile(t *testing.T) {
	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
	}

	e := configs.Model{
		GeneralConfigFolder:    common.DefGSPath,
		SimulatorsConfigFolder: common.DefSCPath,
	}

	a := configs.Configs(&fl)

	s, err := os.Stat(file)

	assert.NotNil(t, err)

	assert.Nil(t, s)

	assert.Equal(t, e, a)

	os.Remove(common.GetThePath(file))

	configsNoFileWithFlags(t)

	os.Remove(common.GetThePath(file))

	updateConfigsNoForce(t)

	os.Remove(common.GetThePath(file))

	updateConfigsForce(t)
}

func configsNoFileWithFlags(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    generalFolder,
		SCFolder:    simulatorFolder,
	}

	e := configs.Model{
		GeneralConfigFolder:    generalFolder,
		SimulatorsConfigFolder: simulatorFolder,
	}

	a := configs.Configs(&fl)

	s, err := os.Stat(file)

	assert.NotNil(t, err)

	assert.Nil(t, s)

	assert.Equal(t, e, a)
}

func updateConfigsNoForce(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
	}

	e := configs.Model{
		GeneralConfigFolder:    common.DefGSPath,
		SimulatorsConfigFolder: common.DefSCPath,
	}

	configs.Configs(&fl)

	fl.GCFolder = generalFolder
	fl.SCFolder = simulatorFolder

	a := configs.Configs(&fl)

	assert.Equal(t, e, a)
}

func updateConfigsForce(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
	}

	e := configs.Model{
		GeneralConfigFolder:    generalFolder,
		SimulatorsConfigFolder: simulatorFolder,
	}

	configs.Configs(&fl)

	fl.GCFolder = generalFolder
	fl.SCFolder = simulatorFolder
	fl.ForceUpdate = true

	a := configs.Configs(&fl)

	assert.Equal(t, e, a)
}
