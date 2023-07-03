//nolint:gocritic,nolintlint
package configs_test

import (
	"os"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/configs"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/stretchr/testify/assert"
)

const (
	file            string = "/configs.json"
	generalFolder   string = "/general"
	simulatorFolder string = "/simulator"
)

func TestConfigs(t *testing.T) {
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

	os.Remove(common.GetThePath(file))

	getConfigsSimPath(t)

	os.Remove(common.GetThePath(file))

	getConfigsSimPathFailRead(t)

	os.Remove(common.GetThePath(file))

	updateConfigsSimPath(t)

	os.Remove(common.GetThePath(file))

	updateConfigsSimPathFailRead(t)
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

func getConfigsSimPath(t *testing.T) {
	t.Helper()

	var fl = flags.Flags{
		ForceUpdate: true,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    simulatorFolder,
	}

	configs.Configs(&fl)

	var e = configs.SimPathModel{
		SimulatorsConfigFolder: simulatorFolder,
	}

	var ok, msg, code, a = configs.GetConfigsSimPath(translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.SimulatorsConfigFolder, a.SimulatorsConfigFolder)
}

func getConfigsSimPathFailRead(t *testing.T) {
	t.Helper()

	var ok, msg, code, a = configs.GetConfigsSimPath(translation.Translation{L: translation.PtPt})

	assert.False(t, ok)

	assert.Equal(t, translation.Translation{L: translation.PtPt}.Get(text.ConfigsErrorRead), msg)

	assert.Equal(t, 500, code)

	assert.Nil(t, a)
}

func updateConfigsSimPath(t *testing.T) {
	t.Helper()

	var fl = flags.Flags{
		ForceUpdate: true,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    simulatorFolder,
	}

	configs.Configs(&fl)

	var e = configs.SimPathModel{
		SimulatorsConfigFolder: "/test",
	}

	var ok, msg, code, a = configs.UpdateConfigsSimPath(&e, translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.SimulatorsConfigFolder, a.SimulatorsConfigFolder)

	ok, msg, code, a = configs.GetConfigsSimPath(translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.SimulatorsConfigFolder, a.SimulatorsConfigFolder)
}

func updateConfigsSimPathFailRead(t *testing.T) {
	t.Helper()

	var e = configs.SimPathModel{
		SimulatorsConfigFolder: "/test",
	}

	var ok, msg, code, a = configs.UpdateConfigsSimPath(&e, translation.Translation{L: translation.PtPt})

	assert.False(t, ok)

	assert.Equal(t, translation.Translation{L: translation.PtPt}.Get(text.ConfigsErrorRead), msg)

	assert.Equal(t, 500, code)

	assert.Nil(t, a)
}
